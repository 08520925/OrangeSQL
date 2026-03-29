import { test, expect } from "@playwright/test";

const API = "http://localhost:5522/api";

async function setupTestData(page: import("@playwright/test").Page) {
  const execSql = async (sql: string) => {
    await page.request.post(`${API}/exec`, { data: { sql } });
  };
  await execSql("CREATE TABLE IF NOT EXISTS phase6_test (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT)");
  await execSql("DELETE FROM phase6_test");
  await execSql("INSERT INTO phase6_test (id, name, email) VALUES (1, 'Alice', 'alice@test.com')");
  await execSql("INSERT INTO phase6_test (id, name, email) VALUES (2, 'Bob', NULL)");
}

async function typeAndExecute(page: import("@playwright/test").Page, sql: string) {
  await page.click(".cm-content");
  await page.keyboard.press("Control+a");
  await page.keyboard.press("Backspace");
  await page.keyboard.type(sql, { delay: 10 });
  await page.waitForTimeout(200);
  await page.keyboard.press("Control+Enter");
  await page.waitForTimeout(1500);
}

test.describe("Phase 6: オートコンプリート", () => {
  test.beforeEach(async ({ page }) => {
    await setupTestData(page);
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("テーブル名の補完候補が表示される", async ({ page }) => {
    // エディタをクリックして Ctrl+Space で補完トリガー
    await page.click(".cm-content");
    await page.keyboard.type("SELECT * FROM pha", { delay: 30 });
    await page.waitForTimeout(500);
    await page.keyboard.press("Control+Space");
    await page.waitForTimeout(500);

    // 補完ポップアップが表示される
    await expect(page.locator(".cm-tooltip-autocomplete")).toBeVisible();
    // phase6_test が候補に含まれる
    await expect(page.locator(".cm-completionLabel")).toContainText(["phase6_test"]);
  });

  test("completions API がテーブルとカラムを返す", async ({ page }) => {
    const res = await page.request.get(`${API}/schema/completions`);
    expect(res.status()).toBe(200);
    const data = await res.json() as { tables: { name: string; columns: { name: string }[] }[] };
    expect(data.tables.length).toBeGreaterThan(0);

    const table = data.tables.find((t: { name: string }) => t.name === "phase6_test");
    expect(table).toBeDefined();
    expect(table?.columns.length).toBe(3);
  });
});

test.describe("Phase 6: 結果エクスポート", () => {
  test.beforeEach(async ({ page }) => {
    await setupTestData(page);
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("CSV / JSON ダウンロードボタンが表示される", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM phase6_test ORDER BY id");

    // エクスポートボタンが表示される
    await expect(page.locator(".export-btn").first()).toBeVisible();
    // CSV, JSON, コピー ボタン
    const buttons = page.locator(".export-btn");
    await expect(buttons).toHaveCount(3);
  });

  test("コピーメニューが開く", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM phase6_test ORDER BY id");

    // コピー ▼ ボタンをクリック
    await page.locator(".export-btn").filter({ hasText: "コピー" }).first().click();

    // メニューが表示される
    await expect(page.locator(".copy-menu")).toBeVisible();
    await expect(page.locator(".copy-item")).toHaveCount(2);
    await expect(page.locator(".copy-item").first()).toContainText("CSV");
    await expect(page.locator(".copy-item").last()).toContainText("Markdown");
  });
});
