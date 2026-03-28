import { ref, onMounted } from "vue";
import { fetchTables } from "../../shared/api";
import type { TableEntry } from "./types";

/**
 * テーブル一覧の取得と更新を管理する composable。
 */
export function useSchema() {
  const tables = ref<TableEntry[]>([]);
  const loading = ref<boolean>(false);

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

  onMounted(() => {
    void refresh();
  });

  return { tables, loading, refresh };
}
