export type SqlEditorApi = {
  getValue: () => string;
  setValue: (sql: string) => void;
  getStatementAtCursor: () => string;
  refreshCompletions: () => Promise<void>;
};
