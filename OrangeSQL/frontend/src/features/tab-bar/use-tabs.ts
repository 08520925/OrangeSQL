import { ref, computed } from "vue";
import type { Tab } from "./types";

let nextId = 1;

function createTab(): Tab {
  const id = String(nextId);
  const title = `Query ${String(nextId)}`;
  nextId++;
  return { id, title, sql: "", result: { kind: "idle" } };
}

/**
 * タブの追加・切り替え・閉じ・状態管理を担う composable。
 */
export function useTabs() {
  const tabs = ref<Tab[]>([createTab()]);
  const activeTabId = ref<string>(tabs.value[0]?.id ?? "1");

  const activeTab = computed<Tab | undefined>(() =>
    tabs.value.find((t) => t.id === activeTabId.value),
  );

  function addTab(): void {
    const tab = createTab();
    tabs.value = [...tabs.value, tab];
    activeTabId.value = tab.id;
  }

  function closeTab(id: string): void {
    if (tabs.value.length <= 1) return;
    const idx = tabs.value.findIndex((t) => t.id === id);
    tabs.value = tabs.value.filter((t) => t.id !== id);
    if (activeTabId.value === id) {
      const newIdx = Math.min(idx, tabs.value.length - 1);
      const newTab = tabs.value[newIdx];
      if (newTab != null) {
        activeTabId.value = newTab.id;
      }
    }
  }

  function switchTab(id: string): void {
    activeTabId.value = id;
  }

  function updateSql(sql: string): void {
    const tab = tabs.value.find((t) => t.id === activeTabId.value);
    if (tab != null) {
      tab.sql = sql;
    }
  }

  function updateResult(result: Tab["result"]): void {
    const tab = tabs.value.find((t) => t.id === activeTabId.value);
    if (tab != null) {
      tab.result = result;
    }
  }

  return { tabs, activeTabId, activeTab, addTab, closeTab, switchTab, updateSql, updateResult };
}
