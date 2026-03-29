import { test, expect } from "@playwright/test";

// テスト用ヘルパー: エディタにSQLを入力して実行
async function typeAndExecute(page: import("@playwright/test").Page, sql: string) {
  // エディタをクリアして入力
  await page.click(".cm-content");
  await page.keyboard.press("Control+a");
  await page.keyboard.press("Backspace");
  await page.keyboard.type(sql, { delay: 10 });
  await page.waitForTimeout(200);
  await page.keyboard.press("Control+Enter");
  await page.waitForTimeout(1500);
}

// テスト用ヘルパー: API 経由でテーブルとデータを準備（UI操作より高速）
async function setupTestData(page: import("@playwright/test").Page) {
  const apiBase = "http://localhost:5522/api";
  const execSql = async (sql: string) => {
    await page.request.post(`${apiBase}/exec`, {
      data: { sql },
    });
  };
  await execSql("CREATE TABLE IF NOT EXISTS test_e2e (id INTEGER PRIMARY KEY, name TEXT, email TEXT)");
  await execSql("DELETE FROM test_e2e");
  await execSql("INSERT INTO test_e2e (name, email) VALUES ('Alice', 'alice@example.com')");
  await execSql("INSERT INTO test_e2e (name, email) VALUES ('Bob', NULL)");
}

test.describe("OrangeSQL ワークベンチ", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("画面が正しく表示される", async ({ page }) => {
    // ヘッダー
    await expect(page.locator(".app-name")).toHaveText("OrangeSQL");
    await expect(page.locator(".dropdown-trigger")).toBeVisible();

    // 実行ボタン
    await expect(page.locator(".execute-btn")).toBeVisible();

    // エディタ
    await expect(page.locator(".cm-editor")).toBeVisible();

    // ステータスバー
    await expect(page.locator(".status-bar")).toContainText("Ready");

    // サイドバー
    await expect(page.locator(".sidebar-title")).toHaveText("テーブル");
  });

  test("SELECT を実行して結果テーブルが表示される", async ({ page }) => {
    await setupTestData(page);
    await typeAndExecute(page, "SELECT * FROM test_e2e ORDER BY id");

    // カラムヘッダ
    const headers = page.locator("th");
    await expect(headers.nth(0)).toContainText("id");
    await expect(headers.nth(1)).toContainText("name");
    await expect(headers.nth(2)).toContainText("email");

    // 行数
    const rows = page.locator("tbody tr");
    await expect(rows).toHaveCount(2);

    // ステータスバー
    await expect(page.locator(".status-bar")).toContainText("2 行取得");
  });

  test("NULL がグレー+イタリックで表示される", async ({ page }) => {
    await setupTestData(page);
    await typeAndExecute(page, "SELECT * FROM test_e2e WHERE name = 'Bob'");

    const nullCell = page.locator(".null-cell");
    await expect(nullCell).toHaveText("NULL");
    await expect(nullCell).toBeVisible();
  });

  test("INSERT を実行して affected rows が表示される", async ({ page }) => {
    await setupTestData(page);
    await typeAndExecute(page, "INSERT INTO test_e2e (name, email) VALUES ('Charlie', 'c@test.com')");

    await expect(page.locator(".status-bar")).toContainText("1 行に影響");
  });

  test("不正な SQL でエラーが表示される", async ({ page }) => {
    await typeAndExecute(page, "SELEC * FROM users");

    await expect(page.locator(".status-bar")).toContainText("エラー");
    await expect(page.locator(".status-bar")).toHaveClass(/status-error/);
  });

  test("サイドバーのテーブルクリックで SELECT 文が挿入される", async ({ page }) => {
    await setupTestData(page);

    // サイドバーを更新（リフレッシュボタン）
    await page.click(".refresh-btn");
    await page.waitForTimeout(1000);

    // test_e2e テーブルをクリック
    const tableItem = page.locator(".table-item", { hasText: "test_e2e" });
    await tableItem.click();
    await page.waitForTimeout(500);

    // エディタに SELECT 文が入っている
    const editorContent = await page.locator(".cm-content").textContent();
    expect(editorContent).toContain("SELECT * FROM test_e2e LIMIT 100");
  });

  test("実行ボタンでも SQL を実行できる", async ({ page }) => {
    await setupTestData(page);

    await page.click(".cm-content");
    await page.keyboard.press("Control+a");
    await page.keyboard.press("Backspace");
    await page.keyboard.type("SELECT count(*) FROM test_e2e", { delay: 10 });
    await page.waitForTimeout(200);

    // 実行ボタンをクリック
    await page.click(".execute-btn");
    await page.waitForTimeout(1500);

    await expect(page.locator(".status-bar")).toContainText("1 行取得");
  });
});
