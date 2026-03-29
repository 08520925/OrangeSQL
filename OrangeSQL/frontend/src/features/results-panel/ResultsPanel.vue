<script setup lang="ts">
import type { ResultState } from "./types";
import type { QueryResult } from "../../shared/types";
import EditableTable from "./EditableTable.vue";

const props = defineProps<{
  state: ResultState;
  lastSql?: string;
}>();

const emit = defineEmits<{
  refreshResult: [data: QueryResult];
  editError: [message: string];
}>();

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
              <td v-for="(cell, j) in row" :key="j" :class="{ 'null-cell': cell === null }">
                {{ cell === null ? 'NULL' : cell }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>
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

.table-wrapper {
  overflow: auto;
  height: 100%;
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
}

tr:hover td {
  background: #2a2d2e;
}

.null-cell {
  color: #666666;
  font-style: italic;
}
</style>
