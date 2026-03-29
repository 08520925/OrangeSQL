/** POST /api/query のレスポンス */
export type QueryResult = {
  columns: string[];
  rows: (string | null)[][];
  rowCount: number;
  executionTimeMs: number;
  truncated?: boolean;
  editableTable?: string;
  pkColumns?: string[];
};

/** POST /api/exec のレスポンス */
export type ExecResult = {
  affectedRows: number;
  executionTimeMs: number;
};

/** POST /api/exec/batch のレスポンス */
export type BatchExecResult = {
  totalAffectedRows: number;
  executionTimeMs: number;
};

/** エラー時のレスポンス（両API共通） */
export type ApiError = {
  error: string;
};
