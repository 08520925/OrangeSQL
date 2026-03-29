/**
 * 結果テーブルを CSV / JSON / Markdown に変換し、ダウンロードまたはクリップボードコピーする。
 */

/** CSV 用: 値をエスケープする。カンマ・改行・ダブルクォートを含む場合はダブルクォートで囲む。 */
function csvEscape(value: string): string {
  if (value.includes(",") || value.includes("\n") || value.includes('"')) {
    return `"${value.replace(/"/g, '""')}"`;
  }
  return value;
}

/** Markdown 用: パイプをエスケープ */
function mdEscape(value: string): string {
  return value.replace(/\|/g, "\\|");
}

/** CSV 文字列を生成する（カラム名ヘッダー付き） */
export function toCSV(columns: string[], rows: (string | null)[][]): string {
  const header = columns.map(csvEscape).join(",");
  const body = rows.map((row) =>
    row.map((cell) => (cell === null ? "" : csvEscape(cell))).join(","),
  );
  return [header, ...body].join("\n");
}

/** JSON 文字列を生成する（オブジェクト配列） */
export function toJSON(columns: string[], rows: (string | null)[][]): string {
  const objects = rows.map((row) => {
    const obj: Record<string, string | null> = {};
    for (let i = 0; i < columns.length; i++) {
      const col = columns[i];
      if (col != null) {
        obj[col] = row[i] ?? null;
      }
    }
    return obj;
  });
  return JSON.stringify(objects, null, 2);
}

/** Markdown テーブル文字列を生成する */
export function toMarkdown(columns: string[], rows: (string | null)[][]): string {
  const header = "| " + columns.map(mdEscape).join(" | ") + " |";
  const separator = "| " + columns.map(() => "---").join(" | ") + " |";
  const body = rows.map(
    (row) =>
      "| " +
      row.map((cell) => (cell === null ? "" : mdEscape(cell))).join(" | ") +
      " |",
  );
  return [header, separator, ...body].join("\n");
}

/** ファイルをダウンロードする */
function downloadFile(content: string, filename: string, mimeType: string): void {
  const blob = new Blob([content], { type: mimeType });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

/** CSV ファイルとしてダウンロード */
export function downloadCSV(columns: string[], rows: (string | null)[][]): void {
  const csv = toCSV(columns, rows);
  downloadFile(csv, "result.csv", "text/csv;charset=utf-8");
}

/** JSON ファイルとしてダウンロード */
export function downloadJSON(columns: string[], rows: (string | null)[][]): void {
  const json = toJSON(columns, rows);
  downloadFile(json, "result.json", "application/json;charset=utf-8");
}

/** クリップボードにコピーする */
async function copyToClipboard(text: string): Promise<void> {
  await navigator.clipboard.writeText(text);
}

/** CSV としてクリップボードにコピー */
export async function copyCSV(columns: string[], rows: (string | null)[][]): Promise<void> {
  await copyToClipboard(toCSV(columns, rows));
}

/** Markdown テーブルとしてクリップボードにコピー */
export async function copyMarkdown(columns: string[], rows: (string | null)[][]): Promise<void> {
  await copyToClipboard(toMarkdown(columns, rows));
}
