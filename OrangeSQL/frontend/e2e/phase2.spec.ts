import { test, expect } from "@playwright/test";

// テスト用ヘルパー: API でテーブルとデータを準備
async function setupTestData(page: import("@playwright/test").Page) {
  const apiBase = "http://localhost:5522/api";
  const execSql = async (sql: string) => {
    await page.request.post(`${apiBase}/exec`, { data: { sql } });
  };
  await execSql("CREATE TABLE IF NOT EXISTS test_p2 (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT)");
  await execSql("DELETE FROM test_p2");
  await execSql("INSERT INTO test_p2 (name, email) VALUES ('Alice', 'a@test.com')");
  await execSql("INSERT INTO test_p2 (name, email) VALUES ('Bob', NULL)");
}

// ヘルパー: エディタにSQLを入力して実行
async function typeAndExecute(page: import("@playwright/test").Page, sql: string) {
  await page.click(".cm-content");
  await page.keyboard.press("Control+a");
  await page.keyboard.press("Backspace");
  await page.keyboard.type(sql, { delay: 10 });
  await page.waitForTimeout(200);
  await page.keyboard.press("Control+Enter");
  await page.waitForTimeout(1500);
}

test.describe("Phase 2: タブ機能", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("新しいタブを追加できる", async ({ page }) => {
    // 初期状態は1タブ
    await expect(page.locator(".tab")).toHaveCount(1);

    // + ボタンでタブ追加
    await page.click(".tab-add");
    await expect(page.locator(".tab")).toHaveCount(2);

    await page.click(".tab-add");
    await expect(page.locator(".tab")).toHaveCount(3);
  });

  test("タブ間でSQLが独立している", async ({ page }) => {
    // Tab1 で SQL を実行
    await typeAndExecute(page, "SELECT 1 AS tab1_result");
    await expect(page.locator(".status-bar")).toContainText("1 行取得");

    // Tab2 を作成して別の SQL を実行
    await page.click(".tab-add");
    await page.waitForTimeout(1000);
    await typeAndExecute(page, "SELECT 2 AS tab2_result");
    await expect(page.locator("th")).toContainText(["tab2_result"]);

    // Tab1 に戻ると Tab1 の結果が表示される
    await page.locator(".tab").first().click();
    await page.waitForTimeout(1000);
    await expect(page.locator("th")).toContainText(["tab1_result"]);
  });

  test("タブを閉じられる（最後の1タブは閉じられない）", async ({ page }) => {
    // タブ追加
    await page.click(".tab-add");
    await expect(page.locator(".tab")).toHaveCount(2);

    // 閉じる
    await page.locator(".tab-close").first().click();
    await expect(page.locator(".tab")).toHaveCount(1);

    // 最後の1タブには × ボタンがない
    await expect(page.locator(".tab-close")).toHaveCount(0);
  });
});

test.describe("Phase 2: カラム情報", () => {
  test.beforeEach(async ({ page }) => {
    await setupTestData(page);
    await page.goto("/");
    await page.waitForTimeout(2000);
    // サイドバーをリフレッシュ
    await page.click(".refresh-btn");
    await page.waitForTimeout(1000);
  });

  test("テーブルを展開するとカラム一覧が表示される", async ({ page }) => {
    // test_p2 の展開ボタンをクリック
    const expandBtn = page.locator(".table-group").filter({ hasText: "test_p2" }).locator(".expand-btn");
    await expandBtn.click();
    await page.waitForTimeout(2000);

    // カラムが表示される
    const columns = page.locator(".table-group").filter({ hasText: "test_p2" }).locator(".column-item");
    await expect(columns).toHaveCount(3, { timeout: 5000 });

    // カラム名を確認
    await expect(columns.nth(0)).toContainText("id");
    await expect(columns.nth(0)).toContainText("INTEGER");
    await expect(columns.nth(1)).toContainText("name");
    await expect(columns.nth(1)).toContainText("NOT NULL");
    await expect(columns.nth(2)).toContainText("email");
  });

  test("展開を折りたためる", async ({ page }) => {
    const tableGroup = page.locator(".table-group").filter({ hasText: "test_p2" });
    const expandBtn = tableGroup.locator(".expand-btn");

    // 展開
    await expandBtn.click();
    await page.waitForTimeout(2000);
    await expect(tableGroup.locator(".column-item")).toHaveCount(3, { timeout: 5000 });

    // 折りたたみ
    await expandBtn.click();
    await page.waitForTimeout(500);
    await expect(tableGroup.locator(".column-item")).toHaveCount(0);
  });
});

test.describe("Phase 2: パネルリサイズ", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("リサイズハンドルが表示される", async ({ page }) => {
    await expect(page.locator(".resize-handle")).toBeVisible();
  });

  test("ドラッグでエディタの高さが変わる", async ({ page }) => {
    const handle = page.locator(".resize-handle");
    const editorBefore = await page.locator(".editor-area").boundingBox();

    // ハンドルを下にドラッグ
    const box = await handle.boundingBox();
    if (box != null && editorBefore != null) {
      await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2);
      await page.mouse.down();
      await page.mouse.move(box.x + box.width / 2, box.y + 100);
      await page.mouse.up();
      await page.waitForTimeout(300);

      const editorAfter = await page.locator(".editor-area").boundingBox();
      if (editorAfter != null) {
        expect(editorAfter.height).toBeGreaterThan(editorBefore.height);
      }
    }
  });
});
