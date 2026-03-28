<script setup lang="ts">
import { ref, watch, nextTick } from "vue";
import TabBar from "./features/tab-bar/TabBar.vue";
import SqlEditor from "./features/sql-editor/SqlEditor.vue";
import ResultsPanel from "./features/results-panel/ResultsPanel.vue";
import SchemaSidebar from "./features/schema-sidebar/SchemaSidebar.vue";
import ResizeHandle from "./features/resize-handle/ResizeHandle.vue";
import { useResults } from "./features/results-panel/use-results";
import { useSchema } from "./features/schema-sidebar/use-schema";
import { useTabs } from "./features/tab-bar/use-tabs";
import { useResize } from "./features/resize-handle/use-resize";
import { fetchHealth } from "./shared/api";

const editorRef = ref<InstanceType<typeof SqlEditor> | null>(null);
const { state, execute } = useResults();
const { tables, loading: schemaLoading, refresh: refreshSchema, toggleExpand, isExpanded, getColumns } = useSchema();
const { tabs, activeTabId, activeTab, addTab, closeTab, switchTab, updateSql, updateResult } = useTabs();
const { editorRatio, onMouseDown } = useResize(".content-area");

const dbName = ref<string>("");

fetchHealth()
  .then((h) => { dbName.value = h.database; })
  .catch(() => { dbName.value = "未接続"; });

// タブ切り替え時にエディタの内容を復元
watch(activeTabId, () => {
  void nextTick(() => {
    const tab = activeTab.value;
    if (tab != null && editorRef.value != null) {
      editorRef.value.setValue(tab.sql);
      state.value = tab.result;
    }
  });
});

// 結果の状態をアクティブタブに同期
watch(state, (s) => {
  updateResult(s);
  if (s.kind === "exec") {
    void refreshSchema();
  }
}, { deep: true });

async function handleExecute(): Promise<void> {
  const sql = editorRef.value?.getValue() ?? "";
  updateSql(sql);
  await execute(sql);
}

function handleSqlChange(): void {
  const sql = editorRef.value?.getValue() ?? "";
  updateSql(sql);
}

function handleSelectTable(tableName: string): void {
  const sql = `SELECT * FROM ${tableName} LIMIT 100`;
  editorRef.value?.setValue(sql);
  updateSql(sql);
}

function handleSwitchTab(id: string): void {
  // 現在のタブのSQLを保存
  const currentSql = editorRef.value?.getValue() ?? "";
  updateSql(currentSql);
  switchTab(id);
}

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
    <TabBar
      :tabs="tabs"
      :active-tab-id="activeTabId"
      @switch-tab="handleSwitchTab"
      @add-tab="addTab"
      @close-tab="closeTab"
    />
    <div class="main-area">
      <aside class="sidebar">
        <SchemaSidebar
          :tables="tables"
          :loading="schemaLoading"
          :is-expanded="isExpanded"
          :get-columns="getColumns"
          @select-table="handleSelectTable"
          @toggle-expand="toggleExpand"
          @refresh="refreshSchema"
        />
      </aside>
      <div class="content-area">
        <div class="editor-area" :style="{ flex: `0 0 ${editorRatio * 100}%` }">
          <SqlEditor ref="editorRef" @execute="handleExecute" @change="handleSqlChange" />
        </div>
        <ResizeHandle @mousedown="onMouseDown" />
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

.app-name { font-weight: 600; font-size: 14px; color: #e0e0e0; }
.db-info { color: #888888; font-size: 12px; }

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
.execute-btn:hover { background: #1a8ad4; }
.execute-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.main-area { display: flex; flex: 1; min-height: 0; }

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

.editor-area { min-height: 100px; overflow: hidden; }

.results-area { flex: 1; min-height: 100px; overflow: auto; }

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
.status-bar.status-error { background: #c72e2e; }
</style>
