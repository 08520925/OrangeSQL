import { ref } from "vue";
import type { QueryResult } from "@/entities/query";
import { executeSQL } from "@/shared/api/wails";

export function useRunSql() {
    const running = ref(false);
    const error = ref("");

    async function runSql(sql: string): Promise<QueryResult | null> {
        error.value = "";
        running.value = true;

        try {
            const trimmed = (sql ?? "").trim();
            if (!trimmed) return { columns: [], rows: [] };

            const res = await executeSQL(trimmed);
            return res;
        } catch (e: any) {
            error.value = e?.message ?? String(e);
            return null;
        } finally {
            running.value = false;
        }
    }

    return { runSql, running, error };
}
