/**
 * テーブル編集から UPDATE / INSERT / DELETE 文を生成するモジュール。
 */

/** シングルクォート内のシングルクォートをエスケープし、値をクォートする。NULL は NULL を返す。 */
function quoteValue(value: string | null): string {
  if (value === null) return "NULL";
  return `'${value.replace(/'/g, "''")}'`;
}

/** WHERE 句を PK カラムで生成する。 */
function buildWhereClause(
  pkColumns: string[],
  columns: string[],
  row: (string | null)[],
): string {
  const conditions = pkColumns.map((pk) => {
    const idx = columns.indexOf(pk);
    const val = idx >= 0 ? (row[idx] ?? null) : null;
    if (val === null) return `${pk} IS NULL`;
    return `${pk} = ${quoteValue(val)}`;
  });
  return conditions.join(" AND ");
}

export type CellEdit = {
  rowIndex: number;
  colIndex: number;
  newValue: string | null;
};

export type NewRow = {
  values: (string | null)[];
};

/** UPDATE 文を生成する。変更セルをグループ化して行ごとに1つの UPDATE にまとめる。 */
export function generateUpdates(
  tableName: string,
  columns: string[],
  pkColumns: string[],
  originalRows: (string | null)[][],
  edits: CellEdit[],
): string[] {
  // 行ごとに変更をグループ化
  const byRow = new Map<number, CellEdit[]>();
  for (const edit of edits) {
    const arr = byRow.get(edit.rowIndex) ?? [];
    arr.push(edit);
    byRow.set(edit.rowIndex, arr);
  }

  const statements: string[] = [];
  for (const [rowIndex, rowEdits] of byRow) {
    const row = originalRows[rowIndex];
    if (row === undefined) continue;
    const setClauses = rowEdits.map(
      (e) => `${columns[e.colIndex]} = ${quoteValue(e.newValue)}`,
    );
    const where = buildWhereClause(pkColumns, columns, row);
    statements.push(
      `UPDATE ${tableName} SET ${setClauses.join(", ")} WHERE ${where}`,
    );
  }
  return statements;
}

/** INSERT 文を生成する。PK カラムの値が空の場合は INSERT から除外する。 */
export function generateInserts(
  tableName: string,
  columns: string[],
  pkColumns: string[],
  newRows: NewRow[],
): string[] {
  return newRows
    .filter((row) => row.values.some((v) => v !== null && v !== ""))
    .map((row) => {
      const pairs: { col: string; val: string }[] = [];
      for (let i = 0; i < columns.length; i++) {
        const col = columns[i] ?? "";
        const val = row.values[i] ?? null;
        // PK が空なら除外（AUTOINCREMENT 対応）
        if (pkColumns.includes(col) && (val === null || val === "")) continue;
        pairs.push({ col: col, val: quoteValue(val) });
      }
      if (pairs.length === 0) return "";
      const colList = pairs.map((p) => p.col).join(", ");
      const valList = pairs.map((p) => p.val).join(", ");
      return `INSERT INTO ${tableName} (${colList}) VALUES (${valList})`;
    })
    .filter((s) => s !== "");
}

/** DELETE 文を生成する。 */
export function generateDeletes(
  tableName: string,
  columns: string[],
  pkColumns: string[],
  originalRows: (string | null)[][],
  deletedRowIndices: Set<number>,
): string[] {
  const statements: string[] = [];
  for (const rowIndex of deletedRowIndices) {
    const row = originalRows[rowIndex];
    if (row === undefined) continue;
    const where = buildWhereClause(pkColumns, columns, row);
    statements.push(`DELETE FROM ${tableName} WHERE ${where}`);
  }
  return statements;
}
