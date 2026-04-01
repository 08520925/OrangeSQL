import type { QueryResult, ExecResult } from "../../shared/types";
import type { ColumnEntry } from "../schema-sidebar/types";

export type TableInfoData = {
  tableName: string;
  columns: ColumnEntry[];
};

export type ResultState =
  | { kind: "idle" }
  | { kind: "loading" }
  | { kind: "query"; data: QueryResult }
  | { kind: "exec"; data: ExecResult }
  | { kind: "error"; message: string }
  | { kind: "tableInfo"; data: TableInfoData };
