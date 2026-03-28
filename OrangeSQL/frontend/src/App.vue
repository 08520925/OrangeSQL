<script setup lang="ts">
import { ref, watch } from "vue";
import SqlEditor from "./features/sql-editor/SqlEditor.vue";
import ResultsPanel from "./features/results-panel/ResultsPanel.vue";
import SchemaSidebar from "./features/schema-sidebar/SchemaSidebar.vue";
import { useResults } from "./features/results-panel/use-results";
import { useSchema } from "./features/schema-sidebar/use-schema";
import { fetchHealth } from "./shared/api";

const editorRef = ref<InstanceType<typeof SqlEditor> | null>(null);
const { state, execute } = useResults();
const { tables, loading: schemaLoading, refresh: refreshSchema } = useSchema();

const dbName = ref<string>("");

// ヘルスチェックで DB 名を取得
fetchHealth()
  .then((h) => { dbName.value = h.database; })
  .catch(() => { dbName.value = "未接続"; });

async function handleExecute(): Promise<void> {
  const sql = editorRef.value?.getValue() ?? "";
  await execute(sql);
}

function handleSelectTable(tableName: string): void {
  editorRef.value?.setValue(`SELECT * FROM ${tableName} LIMIT 100`);
}

// exec 成功後にサイドバーを自動更新
watch(state, (s) => {
  if (s.kind === "exec") {
    void refreshSchema();
  }
});

// ステータスバーのテキスト
function statusText(): string {
  switch (state.value.kind) {
    case "idle": return "Ready";
    case "loading": return "実行中...";
    case "error": return `エラー: ${state.value.message}`;
    case "exec": return `${String(state.value.data.affectedRows)} 行に影響 (${String(state.value.data.executionTimeMs)}ms)`;
    case "query": {
      const q = state.value.data;
      const trunc = q.truncated === true ? " [切り詰め]" : "";
      return `${String(q.rowCount)} 行取得 (${String(q.executionTimeMs)}ms)${trunc}`;
    }
  }
}
</script>

<template>
  <div class="app-layout">
    <header class="header-bar">
      <span class="app-name">OrangeSQL</span>
      <span class="db-info">{{ dbName }}</span>
      <button class="execute-btn" :disabled="state.kind === 'loading'" @click="handleExecute">
        ▶ 実行
      </button>
    </header>
    <div class="main-area">
      <aside class="sidebar">
        <SchemaSidebar
          :tables="tables"
          :loading="schemaLoading"
          @select-table="handleSelectTable"
          @refresh="refreshSchema"
        />
      </aside>
      <div class="content-area">
        <div class="editor-area">
          <SqlEditor ref="editorRef" @execute="handleExecute" />
        </div>
        <div class="results-area">
          <ResultsPanel :state="state" />
        </div>
        <div class="status-bar" :class="{ 'status-error': state.kind === 'error' }">
          <span>{{ statusText() }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
}

.app-layout {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #1e1e1e;
  color: #cccccc;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 13px;
}

.header-bar {
  height: 40px;
  display: flex;
  align-items: center;
  padding: 0 16px;
  gap: 16px;
  background: #2d2d2d;
  border-bottom: 1px solid #404040;
  flex-shrink: 0;
}

.app-name {
  font-weight: 600;
  font-size: 14px;
  color: #e0e0e0;
}

.db-info {
  color: #888888;
  font-size: 12px;
}

.execute-btn {
  margin-left: auto;
  background: #007acc;
  color: #ffffff;
  border: none;
  padding: 4px 14px;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}

.execute-btn:hover {
  background: #1a8ad4;
}

.execute-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.main-area {
  display: flex;
  flex: 1;
  min-height: 0;
}

.sidebar {
  width: 250px;
  background: #252526;
  border-right: 1px solid #404040;
  overflow-y: auto;
  flex-shrink: 0;
}

.content-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.editor-area {
  flex: 1;
  min-height: 0;
  border-bottom: 1px solid #404040;
}

.results-area {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

.status-bar {
  height: 28px;
  display: flex;
  align-items: center;
  padding: 0 12px;
  background: #007acc;
  color: #ffffff;
  font-size: 12px;
  flex-shrink: 0;
}

.status-bar.status-error {
  background: #c72e2e;
}
</style>
