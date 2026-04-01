<script setup lang="ts">
import { ref } from "vue";
import type { ResultState } from "./types";
import type { QueryResult } from "../../shared/types";
import EditableTable from "./EditableTable.vue";
import ExportToolbar from "./ExportToolbar.vue";
import CellDetailPopup from "./CellDetailPopup.vue";

const props = defineProps<{
  state: ResultState;
  lastSql?: string;
}>();

const emit = defineEmits<{
  refreshResult: [data: QueryResult];
  editError: [message: string];
}>();

const popupValue = ref<string | null>(null);

function isEditable(state: ResultState): boolean {
  return (
    state.kind === "query" &&
    state.data.editableTable != null &&
    state.data.editableTable !== "" &&
    state.data.pkColumns != null &&
    state.data.pkColumns.length > 0
  );
}

</script>

<template>
  <div class="results-panel">
    <div v-if="state.kind === 'idle'" class="message">
      Ctrl+Enter でSQLを実行
    </div>

    <div v-else-if="state.kind === 'loading'" class="message">
      <span class="spinner"></span> 実行中...
    </div>

    <div v-else-if="state.kind === 'error'" class="message error">
      {{ state.message }}
    </div>

    <div v-else-if="state.kind === 'exec'" class="message">
      {{ state.data.affectedRows }} 行に影響 ({{ state.data.executionTimeMs }}ms)
    </div>

    <div v-else-if="state.kind === 'tableInfo'" class="table-wrapper">
      <div class="table-info-header">
        <span class="table-info-title">{{ state.data.tableName }}</span>
      </div>
      <table>
        <thead>
          <tr>
            <th></th>
            <th>カラム名</th>
            <th>型</th>
            <th>PK</th>
            <th>NOT NULL</th>
            <th>コメント</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="col in state.data.columns" :key="col.name">
            <td class="pk-icon">{{ col.pk ? '🔑' : '' }}</td>
            <td class="col-name-cell">{{ col.name }}</td>
            <td class="col-type-cell">{{ col.type }}</td>
            <td class="check-cell">{{ col.pk ? '✓' : '' }}</td>
            <td class="check-cell">{{ col.notNull ? '✓' : '' }}</td>
            <td class="comment-cell">{{ col.comment }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <template v-else-if="state.kind === 'query'">
      <!-- 編集可能テーブル -->
      <EditableTable
        v-if="isEditable(state)"
        :data="state.data"
        :original-sql="lastSql ?? ''"
        @refresh="(d: QueryResult) => emit('refreshResult', d)"
        @error="(m: string) => emit('editError', m)"
      />
      <!-- 読み取り専用テーブル -->
      <div v-else class="table-wrapper">
        <div class="readonly-toolbar">
          <ExportToolbar :columns="state.data.columns" :rows="state.data.rows" />
        </div>
        <div v-if="state.data.truncated" class="truncated-notice">
          結果が 10,000 行で切り詰められました
        </div>
        <table>
          <thead>
            <tr>
              <th v-for="col in state.data.columns" :key="col">{{ col }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, i) in state.data.rows" :key="i">
              <td
                v-for="(cell, j) in row"
                :key="j"
                :class="{ 'null-cell': cell === null }"
              >
                {{ cell === null ? 'NULL' : cell }}
                <button
                  v-if="cell != null && cell.length > 100"
                  class="expand-btn"
                  @click="popupValue = cell"
                >…</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <CellDetailPopup
      v-if="popupValue != null"
      :value="popupValue"
      @close="popupValue = null"
    />
  </div>
</template>

<style scoped>
.results-panel {
  height: 100%;
  overflow: auto;
  background: #1e1e1e;
}

.message {
  padding: 20px;
  color: #888888;
}

.message.error {
  color: #f44747;
}

.spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid #555;
  border-top-color: #007acc;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  vertical-align: middle;
  margin-right: 6px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.truncated-notice {
  padding: 6px 12px;
  background: #5a3e00;
  color: #ffcc00;
  font-size: 12px;
}

.readonly-toolbar {
  display: flex;
  align-items: center;
  padding: 4px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
}

.table-wrapper {
  overflow: auto;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.table-wrapper table {
  flex: 1;
}

table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
  font-family: "Consolas", "Courier New", monospace;
}

thead {
  position: sticky;
  top: 0;
  z-index: 1;
}

th {
  background: #2d2d2d;
  color: #cccccc;
  padding: 6px 12px;
  text-align: left;
  border-bottom: 1px solid #404040;
  font-weight: 600;
  white-space: nowrap;
}

td {
  padding: 4px 12px;
  border-bottom: 1px solid #333333;
  color: #cccccc;
  white-space: nowrap;
  max-width: 400px;
  overflow: hidden;
  text-overflow: ellipsis;
  position: relative;
}

.expand-btn {
  position: absolute;
  right: 2px;
  top: 50%;
  transform: translateY(-50%);
  background: #007acc;
  color: #ffffff;
  border: none;
  border-radius: 2px;
  cursor: pointer;
  font-size: 11px;
  padding: 1px 5px;
  line-height: 1;
}
.expand-btn:hover { background: #1a8ad4; }

tr:hover td {
  background: #2a2d2e;
}

.null-cell {
  color: #666666;
  font-style: italic;
}

.table-info-header {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
}
.table-info-title {
  font-weight: 600;
  font-size: 14px;
  color: #e0a030;
}
.pk-icon { width: 24px; text-align: center; font-size: 12px; }
.col-name-cell { color: #cccccc; font-weight: 600; }
.col-type-cell { color: #6a9955; }
.check-cell { text-align: center; color: #4ec9b0; }
.comment-cell { color: #999999; }
</style>
