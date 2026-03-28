import type { QueryResult, ExecResult } from "../../shared/types";

export type ResultState =
  | { kind: "idle" }
  | { kind: "loading" }
  | { kind: "query"; data: QueryResult }
  | { kind: "exec"; data: ExecResult }
  | { kind: "error"; message: string };
