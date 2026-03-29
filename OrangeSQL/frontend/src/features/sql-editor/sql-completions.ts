import {
  type CompletionContext,
  type CompletionResult,
  type Completion,
} from "@codemirror/autocomplete";
import { fetchCompletions } from "../../shared/api";

/** SQL キーワード */
const SQL_KEYWORDS = [
  "SELECT", "FROM", "WHERE", "AND", "OR", "NOT", "IN", "LIKE", "BETWEEN",
  "IS", "NULL", "AS", "ON", "JOIN", "LEFT", "RIGHT", "INNER", "OUTER",
  "CROSS", "FULL", "GROUP", "BY", "ORDER", "ASC", "DESC", "HAVING",
  "LIMIT", "OFFSET", "UNION", "ALL", "DISTINCT", "EXISTS",
  "INSERT", "INTO", "VALUES", "UPDATE", "SET", "DELETE",
  "CREATE", "TABLE", "ALTER", "DROP", "INDEX", "VIEW",
  "PRIMARY", "KEY", "FOREIGN", "REFERENCES", "UNIQUE", "CHECK", "DEFAULT",
  "NOT", "NULL", "AUTO_INCREMENT", "SERIAL",
  "INTEGER", "INT", "BIGINT", "SMALLINT", "TEXT", "VARCHAR", "CHAR",
  "BOOLEAN", "FLOAT", "DOUBLE", "DECIMAL", "DATE", "TIMESTAMP",
  "BEGIN", "COMMIT", "ROLLBACK", "TRANSACTION",
  "CASE", "WHEN", "THEN", "ELSE", "END",
  "COUNT", "SUM", "AVG", "MIN", "MAX", "COALESCE",
  "WITH", "RECURSIVE",
];

type TableSchema = {
  name: string;
  type: string;
  columns: { name: string; type: string }[];
};

let cachedTables: TableSchema[] = [];

/** 補完データを API から取得してキャッシュする */
export async function refreshCompletionData(): Promise<void> {
  try {
    const res = await fetchCompletions();
    cachedTables = res.tables;
  } catch {
    // 失敗時はキャッシュをクリアしない
  }
}

/** カーソル直前のワードの開始位置とテキストを取得 */
function getWordBefore(context: CompletionContext): { from: number; word: string } | null {
  const wordMatch = context.matchBefore(/[\w.]+/);
  if (wordMatch == null) return null;
  return { from: wordMatch.from, word: wordMatch.text };
}

/** テーブル名.カラム名 のドット補完か判定 */
function getDotPrefix(word: string): string | null {
  const dotIdx = word.lastIndexOf(".");
  if (dotIdx < 0) return null;
  return word.substring(0, dotIdx).toLowerCase();
}

/** SQL 補完ソース */
export function sqlCompletionSource(context: CompletionContext): CompletionResult | null {
  const wb = getWordBefore(context);
  if (wb == null || wb.word.length === 0) {
    // 明示的トリガー（Ctrl+Space）のみ対応
    if (!context.explicit) return null;
  }

  const from = wb?.from ?? context.pos;
  const word = wb?.word ?? "";

  // ドット補完: テーブル名.カラム名
  const dotPrefix = getDotPrefix(word);
  if (dotPrefix != null) {
    const table = cachedTables.find((t) => t.name.toLowerCase() === dotPrefix);
    if (table == null) return null;
    const dotPos = from + dotPrefix.length + 1;
    const options: Completion[] = table.columns.map((col) => ({
      label: col.name,
      type: "property",
      detail: col.type,
    }));
    return { from: dotPos, options };
  }

  const options: Completion[] = [];

  // テーブル名
  for (const t of cachedTables) {
    options.push({
      label: t.name,
      type: "class",
      detail: t.type,
    });
  }

  // カラム名（全テーブルから）
  const addedColumns = new Set<string>();
  for (const t of cachedTables) {
    for (const col of t.columns) {
      if (!addedColumns.has(col.name)) {
        addedColumns.add(col.name);
        options.push({
          label: col.name,
          type: "property",
          detail: col.type,
        });
      }
    }
  }

  // SQL キーワード
  for (const kw of SQL_KEYWORDS) {
    options.push({
      label: kw,
      type: "keyword",
    });
  }

  return { from, options, validFor: /^\w*$/ };
}
