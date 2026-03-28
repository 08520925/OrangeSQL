import { ref, reactive, onMounted } from "vue";
import { fetchTables, fetchColumns } from "../../shared/api";
import type { TableEntry, ColumnEntry } from "./types";

/**
 * テーブル一覧・カラム情報の取得と管理を担う composable。
 */
export function useSchema() {
  const tables = ref<TableEntry[]>([]);
  const loading = ref<boolean>(false);
  const expandedColumns = reactive<Record<string, ColumnEntry[]>>({});
  const expandedTables = reactive<Set<string>>(new Set());

  async function refresh(): Promise<void> {
    loading.value = true;
    try {
      const res = await fetchTables();
      tables.value = res.tables;
    } catch {
      tables.value = [];
    } finally {
      loading.value = false;
    }
  }

  async function toggleExpand(tableName: string): Promise<void> {
    if (expandedTables.has(tableName)) {
      expandedTables.delete(tableName);
      return;
    }

    try {
      const res = await fetchColumns(tableName);
      expandedColumns[tableName] = res.columns;
      expandedTables.add(tableName);
    } catch {
      // エラー時は展開しない
    }
  }

  function isExpanded(tableName: string): boolean {
    return expandedTables.has(tableName);
  }

  function getColumns(tableName: string): ColumnEntry[] {
    return expandedColumns[tableName] ?? [];
  }

  onMounted(() => {
    void refresh();
  });

  return { tables, loading, refresh, toggleExpand, isExpanded, getColumns };
}
