# Phase 6: SQL オートコンプリート + 結果エクスポート — 仕様書

## 概要

2つの独立した機能を追加する:

1. **SQL オートコンプリート** — テーブル名・カラム名をエディタ内で補完
2. **結果エクスポート** — CSV / JSON ダウンロード + CSV / Markdown クリップボードコピー

---

## 1. SQL オートコンプリート

### 1.1 動作

- エディタで入力中に補完候補がポップアップ表示される
- テーブル名、カラム名、SQL キーワードを補完
- `FROM`、`JOIN`、`INTO` の直後はテーブル名を優先表示
- テーブル名の後のドット（`users.`）でそのテーブルのカラム名を表示
- Ctrl+Space で手動トリガーも可能

### 1.2 補完データの取得

- `/api/schema/tables` でテーブル一覧を取得
- `/api/schema/columns?table=<name>` で各テーブルのカラム一覧を取得
- 接続切り替え時・スキーマ更新時に補完データをリフレッシュ

### 1.3 補完データの API

新規 API を追加する:

#### `GET /api/schema/completions`

全テーブルとその全カラムを一括取得する。

**Response（HTTP 200）:**
```json
{
  "tables": [
    {
      "name": "users",
      "type": "table",
      "columns": [
        { "name": "id", "type": "integer" },
        { "name": "name", "type": "varchar(100)" }
      ]
    }
  ]
}
```

これにより N+1 リクエストを避ける。

### 1.4 CodeMirror 6 補完拡張

`@codemirror/autocomplete` の `autocompletion()` を使用する。

```typescript
import { autocompletion, type CompletionContext } from "@codemirror/autocomplete";
```

カスタム補完ソースを作成:
- SQL キーワード（SELECT, FROM, WHERE, INSERT, UPDATE, DELETE 等）
- テーブル名（type: "class" アイコン）
- カラム名（type: "property" アイコン）

### 1.5 実装ファイル

```
frontend/src/features/sql-editor/
├── use-sql-editor.ts       # 変更: autocompletion 拡張追加
├── sql-completions.ts      # 新規: 補完データ管理 + 補完ソース

internal/schema/
├── handler.go              # 変更: CompletionsHandler 追加
├── service.go              # 変更: BuildCompletionsResponse 追加

internal/server/
├── router.go               # 変更: GET /api/schema/completions ルート追加
```

---

## 2. 結果エクスポート

### 2.1 エクスポート形式

| 形式 | 対象 | アクション |
|------|------|-----------|
| CSV | ファイル | ダウンロード |
| JSON | ファイル | ダウンロード |
| CSV | クリップボード | コピー |
| Markdown テーブル | クリップボード | コピー |

すべてカラム名（ヘッダー行）を含む。

### 2.2 UI

結果テーブルのヘッダー上にエクスポートツールバーを表示する:

```
┌──────────────────────────────────────────────────┐
│ 100 行取得 (5ms)  [CSV▼] [JSON▼] [コピー▼]      │
├──────────────────────────────────────────────────┤
│  id  │  name    │  email           │             │
```

「コピー▼」ドロップダウン:
- CSV としてコピー
- Markdown テーブルとしてコピー

「CSV▼」「JSON▼」はクリックで即ダウンロード。

### 2.3 CSV 形式

```csv
id,name,email
1,田中太郎,tanaka@example.com
2,佐藤花子,sato@example.com
3,高橋美咲,
```

- カンマ区切り
- 値にカンマ・改行・ダブルクォートが含まれる場合はダブルクォートで囲む
- NULL は空文字

### 2.4 JSON 形式

```json
[
  {"id": "1", "name": "田中太郎", "email": "tanaka@example.com"},
  {"id": "2", "name": "佐藤花子", "email": "sato@example.com"},
  {"id": "3", "name": "高橋美咲", "email": null}
]
```

- 各行がオブジェクト、カラム名がキー
- NULL は JSON null

### 2.5 Markdown テーブル形式

```markdown
| id | name | email |
| --- | --- | --- |
| 1 | 田中太郎 | tanaka@example.com |
| 2 | 佐藤花子 | sato@example.com |
| 3 | 高橋美咲 |  |
```

- パイプ区切り
- NULL は空文字
- 値内のパイプは `\|` にエスケープ

### 2.6 実装ファイル

```
frontend/src/features/results-panel/
├── ResultsPanel.vue         # 変更: エクスポートツールバー表示
├── ExportToolbar.vue        # 新規: エクスポートボタン群
├── result-exporter.ts       # 新規: CSV / JSON / Markdown 変換 + ダウンロード / コピー
```

バックエンド変更なし（フロントエンド完結）。

---

## 3. テスト

### 3.1 バックエンド

| テストファイル | テスト内容 |
|--------------|-----------|
| `schema/handler_test.go` | CompletionsHandler のレスポンス形式 |

### 3.2 E2E

| テスト | 内容 |
|--------|------|
| オートコンプリート | テーブル名入力で補完候補が表示される |
| CSV ダウンロード | ボタンクリックでファイルがダウンロードされる |
| クリップボードコピー | コピーボタンでクリップボードに格納される |

---

## 4. Issue 一覧（Phase 6）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 40 | GET /api/schema/completions API | schema/handler.go, service.go | なし |
| 41 | SQL オートコンプリート（CodeMirror 拡張） | sql-editor/sql-completions.ts | #40 |
| 42 | 結果エクスポート（CSV / JSON / Markdown） | results-panel/result-exporter.ts, ExportToolbar.vue | なし |
| 43 | テスト追加（Phase 6 機能） | handler_test.go, e2e/ | #41, #42 |

---

## 5. Phase 6 完了条件

### サーバー
- [ ] `GET /api/schema/completions` が全テーブル + カラムを返す
- [ ] `go test ./internal/...` で全テストがパスする

### フロントエンド
- [ ] エディタでテーブル名・カラム名が補完される
- [ ] `FROM` の後にテーブル名が優先表示される
- [ ] `table.` でカラム名が表示される
- [ ] CSV / JSON ダウンロードが動作する
- [ ] CSV / Markdown テーブル形式でクリップボードにコピーできる
- [ ] `pnpm build` が成功する
- [ ] `pnpm test:e2e` で全テストがパスする
