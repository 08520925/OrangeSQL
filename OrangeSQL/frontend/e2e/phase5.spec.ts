import { test, expect } from "@playwright/test";

const API = "http://localhost:5522/api";

// テスト用テーブル準備
async function setupEditableTable(page: import("@playwright/test").Page) {
  const execSql = async (sql: string) => {
    await page.request.post(`${API}/exec`, { data: { sql } });
  };
  await execSql("DROP TABLE IF EXISTS edit_test");
  await execSql("CREATE TABLE edit_test (id INTEGER PRIMARY KEY, name TEXT NOT NULL, email TEXT)");
  await execSql("INSERT INTO edit_test (id, name, email) VALUES (1, 'Alice', 'alice@test.com')");
  await execSql("INSERT INTO edit_test (id, name, email) VALUES (2, 'Bob', 'bob@test.com')");
  await execSql("INSERT INTO edit_test (id, name, email) VALUES (3, 'Charlie', NULL)");
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

test.describe("Phase 5: テーブルデータ編集", () => {
  test.beforeEach(async ({ page }) => {
    await setupEditableTable(page);
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("編集可能テーブルが表示される（SELECT * FROM）", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM edit_test ORDER BY id");

    // 編集モードバッジが表示される
    await expect(page.locator(".edit-badge")).toContainText("編集モード: edit_test");

    // PK バッジが表示される
    await expect(page.locator(".pk-badge").first()).toHaveText("PK");

    // 行追加ボタンが表示される
    await expect(page.locator(".btn-add")).toBeVisible();
  });

  test("JOINクエリは読み取り専用", async ({ page }) => {
    await typeAndExecute(page, "SELECT e.* FROM edit_test e JOIN edit_test e2 ON e.id = e2.id");

    // 編集モードバッジが表示されない
    await expect(page.locator(".edit-badge")).not.toBeVisible();
  });

  test("PKなしテーブルは読み取り専用", async ({ page }) => {
    // PKなしテーブルを作成
    await page.request.post(`${API}/exec`, {
      data: { sql: "CREATE TABLE IF NOT EXISTS no_pk (name TEXT, value TEXT)" },
    });
    await page.request.post(`${API}/exec`, {
      data: { sql: "INSERT OR IGNORE INTO no_pk VALUES ('a', 'b')" },
    });

    await typeAndExecute(page, "SELECT * FROM no_pk");
    await expect(page.locator(".edit-badge")).not.toBeVisible();
  });

  test("セルをダブルクリックで編集できる", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM edit_test ORDER BY id");

    // name カラム（2番目）の最初の行をダブルクリック
    const cell = page.locator("tbody tr:first-child td:nth-child(2)");
    await cell.dblclick();

    // input が表示される
    const input = cell.locator("input");
    await expect(input).toBeVisible();

    // 値を変更して Enter
    await input.fill("Alice2");
    await input.press("Enter");

    // セルがハイライトされる
    await expect(cell).toHaveClass(/modified-cell/);

    // 「変更を保存」ボタンが表示される
    await expect(page.locator(".btn-save")).toBeVisible();
  });

  test("セル編集 → 保存 → 反映される", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM edit_test ORDER BY id");

    // Alice の name を変更
    const nameCell = page.locator("tbody tr:first-child td:nth-child(2)");
    await nameCell.dblclick();
    const input = nameCell.locator("input");
    await input.fill("Alice_Updated");
    await input.press("Enter");

    // 保存
    await page.click(".btn-save");
    // 確認ダイアログ
    await expect(page.locator(".dialog-title")).toHaveText("変更の確認");
    await page.click(".dialog .btn-primary");

    await page.waitForTimeout(1500);

    // 更新後のデータを確認
    const updatedCell = page.locator("tbody tr:first-child td:nth-child(2) span");
    await expect(updatedCell).toHaveText("Alice_Updated");
  });

  test("行を削除できる", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM edit_test ORDER BY id");

    // 最初の行の削除ボタンをクリック
    await page.locator("tbody tr:first-child .btn-delete").click();

    // 行が削除マークされる
    await expect(page.locator("tbody tr:first-child")).toHaveClass(/deleted-row/);

    // 保存
    await page.click(".btn-save");
    await page.click(".dialog .btn-primary");
    await page.waitForTimeout(1500);

    // 行が2つになる（Alice が削除された）
    const rows = page.locator("tbody tr:not(.new-row)");
    await expect(rows).toHaveCount(2);
  });

  test("変更を破棄できる", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM edit_test ORDER BY id");

    // セルを編集
    const cell = page.locator("tbody tr:first-child td:nth-child(2)");
    await cell.dblclick();
    const input = cell.locator("input");
    await input.fill("Changed");
    await input.press("Enter");

    // 破棄
    await page.click(".btn-discard");

    // ハイライトが消える
    await expect(cell).not.toHaveClass(/modified-cell/);

    // 保存ボタンが消える
    await expect(page.locator(".btn-save")).not.toBeVisible();
  });

  test("行を追加できる", async ({ page }) => {
    await typeAndExecute(page, "SELECT * FROM edit_test ORDER BY id");

    // 行追加
    await page.click(".btn-add");
    await expect(page.locator(".new-row")).toBeVisible();

    // 新規行に入力（id は自動なのでスキップ、name と email を入力）
    const newRowInputs = page.locator(".new-row .cell-input");
    await newRowInputs.nth(1).fill("Dave");
    await newRowInputs.nth(2).fill("dave@test.com");

    // 保存
    await page.click(".btn-save");
    await page.click(".dialog .btn-primary");
    await page.waitForTimeout(1500);

    // 4行になる
    const rows = page.locator("tbody tr:not(.new-row)");
    await expect(rows).toHaveCount(4);
  });

  test("バッチ API がトランザクションでロールバックする", async ({ page }) => {
    // 不正な SQL を含むバッチを送信
    const res = await page.request.post(`${API}/exec/batch`, {
      data: {
        statements: [
          "UPDATE edit_test SET name = 'SHOULD_NOT_PERSIST' WHERE id = 1",
          "INSERT INTO nonexistent_table VALUES (1)",
        ],
      },
    });

    expect(res.status()).toBe(400);

    // ロールバックされて Alice の名前が変わっていないことを確認
    const queryRes = await page.request.post(`${API}/query`, {
      data: { sql: "SELECT name FROM edit_test WHERE id = 1" },
    });
    const data = await queryRes.json() as { rows: (string | null)[][] };
    expect(data.rows[0]?.[0]).toBe("Alice");
  });
});
