import { ExecuteSQL } from "@/wailsjs/go/main/App";
import { type main } from "@/wailsjs/go/models";

export const executeSQL = async (sql: string): Promise<main.QueryResult> => {
  return await ExecuteSQL(sql);
};
