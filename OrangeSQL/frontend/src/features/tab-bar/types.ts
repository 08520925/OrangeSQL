import type { ResultState } from "../results-panel/types";

export type Tab = {
  id: string;
  title: string;
  sql: string;
  result: ResultState;
};
