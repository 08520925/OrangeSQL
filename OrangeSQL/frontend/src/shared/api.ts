import type { QueryResult, ExecResult, ApiError } from "./types";

const BASE_URL = "/api";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });

  const body: unknown = await res.json();

  if (!res.ok) {
    const err = body as ApiError;
    throw new Error(err.error ?? `HTTP ${String(res.status)}`);
  }

  return body as T;
}

/** SELECT 系 SQL を実行する */
export function querySQL(sql: string): Promise<QueryResult> {
  return request<QueryResult>("/query", {
    method: "POST",
    body: JSON.stringify({ sql }),
  });
}

/** INSERT/UPDATE/DELETE/DDL を実行する */
export function execSQL(sql: string): Promise<ExecResult> {
  return request<ExecResult>("/exec", {
    method: "POST",
    body: JSON.stringify({ sql }),
  });
}

/** テーブル一覧を取得する */
export function fetchTables(): Promise<{ tables: { name: string; type: string }[] }> {
  return request("/schema/tables");
}

/** カラム情報を取得する */
export function fetchColumns(table: string): Promise<{ columns: { name: string; type: string; pk: boolean; notNull: boolean }[] }> {
  return request(`/schema/columns?table=${encodeURIComponent(table)}`);
}

/** ヘルスチェック */
export function fetchHealth(): Promise<{ status: string; database: string }> {
  return request("/health");
}
