import { test, expect } from "@playwright/test";

test.describe("Phase 4: マルチDB対応 UI", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
    await page.waitForTimeout(2000);
  });

  test("接続ダイアログでドライバを切り替えるとフォームが変わる", async ({ page }) => {
    // ドロップダウンを開いて「新規接続」をクリック
    await page.click(".dropdown-trigger");
    await page.click("text=+ 新規接続");

    // デフォルトは SQLite: ファイルパスフィールドが表示される
    await expect(page.locator("text=ファイルパス")).toBeVisible();
    await expect(page.locator("text=ホスト")).not.toBeVisible();

    // PostgreSQL に切り替え
    await page.selectOption(".dialog select:not([disabled])", "postgres");
    await expect(page.locator("text=ホスト")).toBeVisible();
    await expect(page.locator("text=ポート")).toBeVisible();
    await expect(page.locator("text=ユーザー")).toBeVisible();
    await expect(page.locator("text=パスワード")).toBeVisible();
    await expect(page.locator("text=データベース")).toBeVisible();
    await expect(page.locator("text=SSL モード")).toBeVisible();
    await expect(page.locator("text=ファイルパス")).not.toBeVisible();

    // MySQL に切り替え: SSL モードは非表示
    await page.selectOption(".dialog select:not([disabled])", "mysql");
    await expect(page.locator("text=ホスト")).toBeVisible();
    await expect(page.locator("text=SSL モード")).not.toBeVisible();

    // SQL Server に切り替え
    await page.selectOption(".dialog select:not([disabled])", "sqlserver");
    await expect(page.locator("text=ホスト")).toBeVisible();
    await expect(page.locator("text=SSL モード")).not.toBeVisible();

    // SQLite に戻す
    await page.selectOption(".dialog select:not([disabled])", "sqlite");
    await expect(page.locator("text=ファイルパス")).toBeVisible();
    await expect(page.locator("text=ホスト")).not.toBeVisible();
  });

  test("SQLite のバリデーション: 接続名とファイルパスが必須", async ({ page }) => {
    await page.click(".dropdown-trigger");
    await page.click("text=+ 新規接続");

    // 何も入力せずに接続ボタンをクリック
    await page.click(".dialog .btn-primary");
    await expect(page.locator(".error")).toBeVisible();

    // 接続名のみ入力
    await page.fill('input[placeholder="例: 開発DB"]', "テストDB");
    await page.click(".dialog .btn-primary");
    await expect(page.locator(".error")).toBeVisible();
  });

  test("TCP 系のバリデーション: ホストとデータベースが必須", async ({ page }) => {
    await page.click(".dropdown-trigger");
    await page.click("text=+ 新規接続");

    // PostgreSQL を選択
    await page.selectOption(".dialog select:not([disabled])", "postgres");
    await page.fill('input[placeholder="例: 開発DB"]', "テストPG");

    // ホストをクリアして接続（ホスト未入力）
    await page.fill('input[placeholder="localhost"]', "");
    await page.click(".dialog .btn-primary");
    await expect(page.locator(".error")).toBeVisible();

    // ホストのみ入力（DB名なし）
    await page.fill('input[placeholder="localhost"]', "localhost");
    await page.click(".dialog .btn-primary");
    await expect(page.locator(".error")).toBeVisible();
  });

  test("接続管理ダイアログにドライバ名が表示される", async ({ page }) => {
    await page.click(".dropdown-trigger");
    await page.click("text=⚙ 接続管理");

    // 既存プロファイル（SQLite）のドライバ名が表示される
    await expect(page.locator(".profile-driver").first()).toBeVisible();
    const driverText = await page.locator(".profile-driver").first().textContent();
    expect(driverText).toBe("sqlite");
  });

  test("SQLite プロファイルを新規作成できる", async ({ page }) => {
    await page.click(".dropdown-trigger");
    await page.click("text=+ 新規接続");

    await page.fill('input[placeholder="例: 開発DB"]', "E2E テストDB");
    await page.fill('input[placeholder="例: C:/work/dev.db"]', ":memory:");
    await page.click(".dialog .btn-primary");

    // ダイアログが閉じる
    await page.waitForTimeout(1500);
    await expect(page.locator(".dialog-backdrop")).not.toBeVisible();

    // ドロップダウンに新しいプロファイルが追加される（ヘッダーに表示される）
    await expect(page.locator(".dropdown-trigger")).toContainText("E2E テストDB");
  });
});
