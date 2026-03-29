import type { QueryResult, ExecResult, BatchExecResult, ApiError } from "./types";

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

/** 複数 SQL をトランザクション内で一括実行する */
export function execBatch(statements: string[]): Promise<BatchExecResult> {
  return request<BatchExecResult>("/exec/batch", {
    method: "POST",
    body: JSON.stringify({ statements }),
  });
}

/** 補完用: 全テーブル + カラムを一括取得 */
export function fetchCompletions(): Promise<{ tables: { name: string; type: string; columns: { name: string; type: string }[] }[] }> {
  return request("/schema/completions");
}

/** ヘルスチェック */
export function fetchHealth(): Promise<{ status: string; database: string }> {
  return request("/health");
}

/** プロファイル一覧を取得する */
export function fetchProfiles(): Promise<{ profiles: { id: string; name: string; driver: string; path?: string; host?: string; port?: number; user?: string; password?: string; dbName?: string; sslMode?: string; createdAt: string }[]; activeId: string }> {
  return request("/profiles");
}

/** プロファイルを作成する */
export function createProfile(data: { name: string; driver: string; path?: string; host?: string; port?: number; user?: string; password?: string; dbName?: string; sslMode?: string }): Promise<unknown> {
  return request("/profiles", { method: "POST", body: JSON.stringify(data) });
}

/** プロファイルを更新する */
export function updateProfile(id: string, data: { name?: string; path?: string; host?: string; port?: number; user?: string; password?: string; dbName?: string; sslMode?: string }): Promise<unknown> {
  return request(`/profiles/${id}`, { method: "PUT", body: JSON.stringify(data) });
}

/** プロファイルを削除する */
export function deleteProfile(id: string): Promise<unknown> {
  return request(`/profiles/${id}`, { method: "DELETE" });
}

/** プロファイルに接続する */
export function connectProfile(id: string): Promise<{ connected: boolean; database: string }> {
  return request(`/profiles/${id}/connect`, { method: "POST" });
}
