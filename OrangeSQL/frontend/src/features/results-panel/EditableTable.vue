<script setup lang="ts">
import { ref } from "vue";
import type { QueryResult } from "../../shared/types";
import { execBatch, querySQL } from "../../shared/api";
import EditableCell from "./EditableCell.vue";
import ExportToolbar from "./ExportToolbar.vue";
import CellDetailPopup from "./CellDetailPopup.vue";
import { useTableEdit } from "./use-table-edit";

const props = defineProps<{
  data: QueryResult;
  originalSql: string;
}>();

const emit = defineEmits<{
  refresh: [data: QueryResult];
  error: [message: string];
}>();

const tableName = props.data.editableTable ?? "";
const columns = props.data.columns;
const pkColumns = props.data.pkColumns ?? [];

const {
  edits,
  newRows,
  deletedRows,
  hasChanges,
  updateCell,
  addNewRow,
  updateNewRowCell,
  toggleDeleteRow,
  discardAll,
  generateAllSQL,
  changeSummary,
} = useTableEdit(tableName, columns, pkColumns, props.data.rows);

const showConfirm = ref<boolean>(false);
const saving = ref<boolean>(false);
const popupValue = ref<string | null>(null);

function isPKColumn(colIndex: number): boolean {
  return pkColumns.includes(columns[colIndex] ?? "");
}

function isModified(rowIndex: number, colIndex: number): boolean {
  return edits.value.some(
    (e) => e.rowIndex === rowIndex && e.colIndex === colIndex,
  );
}

function getCellValue(rowIndex: number, colIndex: number): string | null {
  const edit = edits.value.find(
    (e) => e.rowIndex === rowIndex && e.colIndex === colIndex,
  );
  if (edit != null) return edit.newValue;
  return props.data.rows[rowIndex]?.[colIndex] ?? null;
}

function handleSetNull(rowIndex: number, colIndex: number): void {
  updateCell(rowIndex, colIndex, null);
}

function handleNewRowSetNull(newRowIndex: number, colIndex: number): void {
  updateNewRowCell(newRowIndex, colIndex, null);
}

async function handleSave(): Promise<void> {
  const statements = generateAllSQL();
  if (statements.length === 0) return;

  saving.value = true;
  showConfirm.value = false;

  try {
    await execBatch(statements);
    // 元の SELECT を再実行して結果を更新
    const refreshed = await querySQL(props.originalSql);
    discardAll();
    emit("refresh", refreshed);
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : String(e);
    emit("error", message);
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <div class="editable-table-wrapper">
    <div class="toolbar">
      <span class="edit-badge">編集モード: {{ tableName }}</span>
      <ExportToolbar :columns="columns" :rows="data.rows" />
      <button class="btn-add" @click="addNewRow">+ 行を追加</button>
      <template v-if="hasChanges">
        <button class="btn-save" :disabled="saving" @click="showConfirm = true">
          変更を保存
        </button>
        <button class="btn-discard" :disabled="saving" @click="discardAll">
          変更を破棄
        </button>
      </template>
    </div>

    <div class="table-scroll">
      <table>
        <thead>
          <tr>
            <th v-for="col in columns" :key="col" :class="{ 'pk-header': pkColumns.includes(col) }">
              {{ col }}
              <span v-if="pkColumns.includes(col)" class="pk-badge">PK</span>
            </th>
            <th class="action-col"></th>
          </tr>
        </thead>
        <tbody>
          <!-- 既存行 -->
          <tr
            v-for="(row, ri) in data.rows"
            :key="ri"
            :class="{ 'deleted-row': deletedRows.has(ri) }"
          >
            <EditableCell
              v-for="(_cell, ci) in row"
              :key="ci"
              :value="getCellValue(ri, ci)"
              :editable="!isPKColumn(ci) && !deletedRows.has(ri)"
              :modified="isModified(ri, ci)"
              @update="(v: string | null) => updateCell(ri, ci, v)"
              @set-null="handleSetNull(ri, ci)"
              @show-detail="(v: string) => popupValue = v"
            />
            <td class="action-col">
              <button
                class="btn-delete"
                :class="{ 'btn-undo-delete': deletedRows.has(ri) }"
                @click="toggleDeleteRow(ri)"
              >
                {{ deletedRows.has(ri) ? '↩' : '🗑' }}
              </button>
            </td>
          </tr>
          <!-- 新規行 -->
          <tr v-for="(newRow, nri) in newRows" :key="'new-' + String(nri)" class="new-row">
            <td
              v-for="(_col, ci) in columns"
              :key="ci"
              :class="{ 'pk-cell': isPKColumn(ci), 'null-new-cell': newRow.values[ci] === null }"
            >
              <div class="new-row-cell">
                <input
                  class="cell-input"
                  :class="{ 'null-input': newRow.values[ci] === null }"
                  :placeholder="isPKColumn(ci) ? '(自動)' : newRow.values[ci] === null ? 'NULL' : ''"
                  :value="newRow.values[ci] === null ? '' : (newRow.values[ci] ?? '')"
                  @input="(e: Event) => updateNewRowCell(nri, ci, (e.target as HTMLInputElement).value || null)"
                  @keydown.ctrl.shift.n.prevent="handleNewRowSetNull(nri, ci)"
                />
                <button
                  v-if="!isPKColumn(ci)"
                  class="null-btn-new"
                  :class="{ 'null-active': newRow.values[ci] === null }"
                  title="NULLを設定"
                  @click="handleNewRowSetNull(nri, ci)"
                >NULL</button>
              </div>
            </td>
            <td class="action-col"></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 確認ダイアログ -->
    <div v-if="showConfirm" class="dialog-backdrop" @click.self="showConfirm = false">
      <div class="dialog">
        <div class="dialog-title">変更の確認</div>
        <div class="dialog-body">
          <p>以下の変更を実行します:</p>
          <ul>
            <li v-if="changeSummary().updates > 0">UPDATE {{ changeSummary().updates }}件</li>
            <li v-if="changeSummary().inserts > 0">INSERT {{ changeSummary().inserts }}件</li>
            <li v-if="changeSummary().deletes > 0">DELETE {{ changeSummary().deletes }}件</li>
          </ul>
          <p class="dialog-note">トランザクション内で実行されます（全成功 or 全ロールバック）</p>
        </div>
        <div class="dialog-actions">
          <button class="btn btn-cancel" @click="showConfirm = false">キャンセル</button>
          <button class="btn btn-primary" @click="handleSave">実行</button>
        </div>
      </div>
    </div>

    <CellDetailPopup
      v-if="popupValue != null"
      :value="popupValue"
      @close="popupValue = null"
    />
  </div>
</template>

<style scoped>
.editable-table-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
  flex-shrink: 0;
}

.edit-badge {
  font-size: 12px;
  color: #e0a030;
  font-weight: 600;
}

.btn-add, .btn-save, .btn-discard {
  padding: 3px 10px;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 11px;
}
.btn-add { background: #3a3a3a; color: #cccccc; }
.btn-add:hover { background: #444444; }
.btn-save { background: #007acc; color: #ffffff; }
.btn-save:hover { background: #1a8ad4; }
.btn-save:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-discard { background: #3a3a3a; color: #cccccc; }
.btn-discard:hover { background: #444444; }

.table-scroll {
  flex: 1;
  overflow: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  font-family: "Consolas", "Courier New", monospace;
}

thead { position: sticky; top: 0; z-index: 1; }

th {
  background: #2d2d2d;
  color: #cccccc;
  padding: 6px 12px;
  text-align: left;
  border-bottom: 1px solid #404040;
  font-weight: 600;
  white-space: nowrap;
}
.pk-header { color: #e0a030; }
.pk-badge {
  font-size: 9px;
  background: #5a4000;
  color: #e0a030;
  padding: 1px 4px;
  border-radius: 2px;
  margin-left: 4px;
}

.action-col {
  width: 40px;
  text-align: center;
  padding: 4px;
}

.btn-delete {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 14px;
  padding: 2px 4px;
  border-radius: 3px;
}
.btn-delete { color: #cc6666; }
.btn-delete:hover { background: #3a3a3a; }
.btn-undo-delete { color: #4ec9b0; }

.deleted-row td {
  text-decoration: line-through;
  background: #4a1a1a;
  color: #cc6666;
}

.new-row td {
  background: #1a2a1a;
}

.pk-cell {
  background: #2a2a2a;
  color: #888888;
}

.cell-input {
  width: 100%;
  padding: 2px 4px;
  background: transparent;
  border: 1px solid #444444;
  color: #cccccc;
  font-size: 13px;
  font-family: "Consolas", "Courier New", monospace;
  box-sizing: border-box;
  outline: none;
}
.cell-input:focus { border-color: #007acc; }
.cell-input.null-input { color: #666666; font-style: italic; }

.new-row-cell {
  display: flex;
  align-items: center;
  gap: 2px;
}

.null-btn-new {
  padding: 1px 4px;
  background: #3a3a3a;
  color: #888888;
  border: 1px solid #555555;
  border-radius: 2px;
  cursor: pointer;
  font-size: 10px;
  font-weight: 600;
  line-height: 1.2;
  white-space: nowrap;
  flex-shrink: 0;
}
.null-btn-new:hover { background: #5a4000; color: #e0a030; }
.null-btn-new.null-active { background: #5a4000; color: #e0a030; border-color: #7a6020; }

tr:hover td { background: #2a2d2e; }
.deleted-row:hover td { background: #552222; }
.new-row:hover td { background: #1a2a1a; }

/* Dialog */
.dialog-backdrop {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex; align-items: center; justify-content: center;
  z-index: 200;
}
.dialog {
  background: #2d2d2d; border: 1px solid #555555; border-radius: 6px;
  padding: 20px; width: 350px; max-width: 90vw;
}
.dialog-title { font-size: 14px; font-weight: 600; color: #e0e0e0; margin-bottom: 12px; }
.dialog-body { font-size: 13px; color: #cccccc; }
.dialog-body ul { margin: 8px 0; padding-left: 20px; }
.dialog-body li { margin: 4px 0; }
.dialog-note { font-size: 11px; color: #888888; margin-top: 8px; }
.dialog-actions { display: flex; justify-content: flex-end; gap: 8px; margin-top: 16px; }
.btn { padding: 6px 16px; border: none; border-radius: 3px; cursor: pointer; font-size: 12px; }
.btn-cancel { background: #3a3a3a; color: #cccccc; }
.btn-cancel:hover { background: #444444; }
.btn-primary { background: #007acc; color: #ffffff; }
.btn-primary:hover { background: #1a8ad4; }
</style>
