<script setup lang="ts">
import { ref, watch, nextTick } from "vue";
import type { QueryResult } from "./shared/types";
import { fetchColumns } from "./shared/api";
import TabBar from "./features/tab-bar/TabBar.vue";
import SqlEditor from "./features/sql-editor/SqlEditor.vue";
import ResultsPanel from "./features/results-panel/ResultsPanel.vue";
import SchemaSidebar from "./features/schema-sidebar/SchemaSidebar.vue";
import ResizeHandle from "./features/resize-handle/ResizeHandle.vue";
import ConnectionDropdown from "./features/connection/ConnectionDropdown.vue";
import ConnectionDialog from "./features/connection/ConnectionDialog.vue";
import ConnectionManagerVue from "./features/connection/ConnectionManager.vue";
import { useResults } from "./features/results-panel/use-results";
import { useSchema } from "./features/schema-sidebar/use-schema";
import { useTabs } from "./features/tab-bar/use-tabs";
import { useResize } from "./features/resize-handle/use-resize";
import { useConnection } from "./features/connection/use-connection";
import type { ConnectionProfile, CreateProfileRequest } from "./features/connection/types";

const editorRef = ref<InstanceType<typeof SqlEditor> | null>(null);
const { state, execute } = useResults();
const { tables, loading: schemaLoading, refresh: refreshSchema } = useSchema();
const { tabs, activeTabId, activeTab, addTab, closeTab, switchTab, updateSql, updateResult } = useTabs();
const { editorRatio, onMouseDown } = useResize(".content-area");
const { profiles, activeId, connect, create, update, remove } = useConnection();

// 最後に実行した SQL（編集モードの再取得用）
const lastExecutedSql = ref<string>("");

// ダイアログ状態
const showNewDialog = ref<boolean>(false);
const showManagerDialog = ref<boolean>(false);
const editingProfile = ref<ConnectionProfile | undefined>(undefined);

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
    editorRef.value?.refreshCompletions();
  }
}, { deep: true });

async function handleExecute(): Promise<void> {
  const fullSql = editorRef.value?.getValue() ?? "";
  updateSql(fullSql);
  const sql = editorRef.value?.getStatementAtCursor() ?? "";
  lastExecutedSql.value = sql;
  await execute(sql);
}

function handleRefreshResult(data: QueryResult): void {
  state.value = { kind: "query", data };
}

function handleEditError(message: string): void {
  state.value = { kind: "error", message };
}

function handleSelectTable(tableName: string): void {
  const activeProfile = profiles.value.find((p) => p.id === activeId.value);
  const sql = activeProfile?.driver === "sqlserver"
    ? `SELECT TOP 100 * FROM ${tableName};`
    : `SELECT * FROM ${tableName} LIMIT 100;`;
  // テーブル情報タブだった場合、結果をリセットしてエディタを表示させる
  if (state.value.kind === "tableInfo") {
    state.value = { kind: "idle" };
  }
  void nextTick(() => {
    editorRef.value?.setValue(sql);
    updateSql(sql);
  });
}

async function handleShowTableInfo(tableName: string): Promise<void> {
  addTab();
  const tab = activeTab.value;
  if (tab != null) {
    tab.title = `${tableName}(情報)`;
  }
  try {
    const res = await fetchColumns(tableName);
    state.value = { kind: "tableInfo", data: { tableName, columns: res.columns } };
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : String(e);
    state.value = { kind: "error", message };
  }
}

function handleSwitchTab(id: string): void {
  const currentSql = editorRef.value?.getValue() ?? "";
  updateSql(currentSql);
  switchTab(id);
}

async function handleConnect(id: string): Promise<void> {
  await connect(id);
  void refreshSchema();
  editorRef.value?.refreshCompletions();
  // 全タブの結果をリセット
  for (const tab of tabs.value) {
    tab.result = { kind: "idle" };
  }
  state.value = { kind: "idle" };
}

async function handleCreateProfile(req: CreateProfileRequest): Promise<void> {
  await create(req);
  showNewDialog.value = false;
  // 作成後に最新プロファイルに接続
  const latest = profiles.value[profiles.value.length - 1];
  if (latest != null) {
    await handleConnect(latest.id);
  }
}

async function handleUpdateProfile(req: CreateProfileRequest): Promise<void> {
  if (editingProfile.value != null) {
    await update(editingProfile.value.id, req);
    editingProfile.value = undefined;
  }
}

async function handleDeleteProfile(id: string): Promise<void> {
  await remove(id);
}

function handleEdit(p: ConnectionProfile): void {
  showManagerDialog.value = false;
  editingProfile.value = p;
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
    case "tableInfo":
      return `テーブル情報: ${state.value.data.tableName} (${String(state.value.data.columns.length)} カラム)`;
  }
}
</script>

<template>
  <div class="app-layout">
    <header class="header-bar">
      <span class="app-name">OrangeSQL</span>
      <ConnectionDropdown
        :profiles="profiles"
        :active-id="activeId"
        @connect="handleConnect"
        @open-new="showNewDialog = true"
        @open-manager="showManagerDialog = true"
      />
      <button class="format-btn" title="SQL整形 (Ctrl+Q)" @click="editorRef?.formatSql()">
        整形
      </button>
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
          @select-table="handleSelectTable"
          @show-table-info="handleShowTableInfo"
          @refresh="refreshSchema"
        />
      </aside>
      <div class="content-area">
        <template v-if="state.kind !== 'tableInfo'">
          <div class="editor-area" :style="{ flex: `0 0 ${editorRatio * 100}%` }">
            <SqlEditor ref="editorRef" @execute="handleExecute" />
          </div>
          <ResizeHandle @mousedown="onMouseDown" />
        </template>
        <div class="results-area">
          <ResultsPanel
            :state="state"
            :last-sql="lastExecutedSql"
            @refresh-result="handleRefreshResult"
            @edit-error="handleEditError"
          />
        </div>
        <div class="status-bar" :class="{ 'status-error': state.kind === 'error' }">
          <span>{{ statusText() }}</span>
        </div>
      </div>
    </div>

    <!-- ダイアログ -->
    <ConnectionDialog
      v-if="showNewDialog"
      @save="handleCreateProfile"
      @close="showNewDialog = false"
    />
    <ConnectionDialog
      v-if="editingProfile != null"
      :edit-profile="editingProfile"
      @save="handleUpdateProfile"
      @close="editingProfile = undefined"
    />
    <ConnectionManagerVue
      v-if="showManagerDialog"
      :profiles="profiles"
      :active-id="activeId"
      @edit="handleEdit"
      @delete="handleDeleteProfile"
      @close="showManagerDialog = false"
    />
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
html, body, #app { height: 100%; width: 100%; overflow: hidden; }

.app-layout {
  display: flex; flex-direction: column; height: 100%;
  background: #1e1e1e; color: #cccccc;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 13px;
}

.header-bar {
  height: 40px; display: flex; align-items: center;
  padding: 0 16px; gap: 16px;
  background: #2d2d2d; border-bottom: 1px solid #404040; flex-shrink: 0;
}
.app-name { font-weight: 600; font-size: 14px; color: #e0e0e0; }

.execute-btn {
  margin-left: auto; background: #007acc; color: #ffffff;
  border: none; padding: 4px 14px; border-radius: 3px;
  cursor: pointer; font-size: 12px;
}
.execute-btn:hover { background: #1a8ad4; }
.execute-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.format-btn {
  background: #3c3c3c; color: #cccccc;
  border: 1px solid #555555; padding: 4px 12px; border-radius: 3px;
  cursor: pointer; font-size: 12px;
}
.format-btn:hover { background: #4a4a4a; }

.main-area { display: flex; flex: 1; min-height: 0; }

.sidebar {
  width: 250px; background: #252526;
  border-right: 1px solid #404040; overflow-y: auto; flex-shrink: 0;
}

.content-area { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.editor-area { min-height: 100px; overflow: hidden; }
.results-area { flex: 1; min-height: 100px; overflow: auto; }

.status-bar {
  height: 28px; display: flex; align-items: center;
  padding: 0 12px; background: #007acc; color: #ffffff;
  font-size: 12px; flex-shrink: 0;
}
.status-bar.status-error { background: #c72e2e; }
</style>
