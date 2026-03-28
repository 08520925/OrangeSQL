import { defineConfig } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  timeout: 30000,
  use: {
    baseURL: "http://localhost:5173",
    headless: true,
    viewport: { width: 1280, height: 720 },
  },
  // テスト前にサーバーが起動済みであることを前提とする
  // バックエンド: go run . (ポート5522)
  // フロントエンド: pnpm dev (ポート5173)
});
