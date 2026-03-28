import { ref } from "vue";
import { querySQL, execSQL } from "../../shared/api";
import type { ResultState } from "./types";

const QUERY_KEYWORDS = new Set(["SELECT", "SHOW", "EXPLAIN", "PRAGMA", "WITH"]);

function getFirstKeyword(sql: string): string {
  const trimmed = sql.trimStart();
  const firstWord = trimmed.split(/\s+/)[0];
  return firstWord?.toUpperCase() ?? "";
}

/**
 * SQL 実行と結果管理の composable。
 */
export function useResults() {
  const state = ref<ResultState>({ kind: "idle" });

  async function execute(sql: string): Promise<void> {
    if (state.value.kind === "loading") return;

    const trimmed = sql.trim();
    if (trimmed === "") return;

    state.value = { kind: "loading" };

    try {
      const keyword = getFirstKeyword(trimmed);
      if (QUERY_KEYWORDS.has(keyword)) {
        const data = await querySQL(trimmed);
        state.value = { kind: "query", data };
      } else {
        const data = await execSQL(trimmed);
        state.value = { kind: "exec", data };
      }
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : String(e);
      state.value = { kind: "error", message };
    }
  }

  return { state, execute };
}
