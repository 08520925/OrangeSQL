/** POST /api/query のレスポンス */
export type QueryResult = {
  columns: string[];
  rows: (string | null)[][];
  rowCount: number;
  executionTimeMs: number;
  truncated?: boolean;
};

/** POST /api/exec のレスポンス */
export type ExecResult = {
  affectedRows: number;
  executionTimeMs: number;
};

/** エラー時のレスポンス（両API共通） */
export type ApiError = {
  error: string;
};
