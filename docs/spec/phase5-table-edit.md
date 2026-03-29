# Phase 5: テーブルデータ編集 — 仕様書

## 概要

結果グリッド上で直接データを編集し、UPDATE / INSERT / DELETE 文を自動生成して実行する。
PK を持つテーブルが対象。PK がないテーブルは編集不可とする。

---

## 1. 編集モードの条件

結果グリッドが「編集可能」になるのは、以下の条件をすべて満たすとき:

1. **単一テーブルの SELECT**（例: `SELECT * FROM users LIMIT 100`）
2. そのテーブルに **PK が存在する**
3. **PK カラムが結果に含まれている**

条件を満たさない場合は、従来通り読み取り専用で表示する。

### 1.1 テーブル名の判定

バックエンドの `/api/query` レスポンスに `editableTable` フィールドを追加する。
バックエンドが SQL を解析し、単一テーブルの SELECT であればテーブル名を返す。

```json
{
  "columns": ["id", "name", "email"],
  "rows": [...],
  "rowCount": 100,
  "executionTimeMs": 5,
  "editableTable": "users",
  "pkColumns": ["id"]
}
```

- `editableTable` が空文字 or 省略 → 読み取り専用
- `pkColumns` が空配列 or 省略 → 読み取り専用

### 1.2 単一テーブル SELECT の判定ロジック

簡易パーサーで判定する（複雑な SQL パーサーは導入しない）:

- `SELECT ... FROM <tableName>` の形式で、JOIN / サブクエリ / UNION がない
- テーブル名を正規表現で抽出: `(?i)FROM\s+([a-zA-Z_]\w*)(?:\s|$|;)`
- 抽出できない場合は editableTable を空にする

PK 情報は `db.Columns(tableName)` で取得し、PK カラムが結果の columns に含まれているか確認する。

---

## 2. UI

### 2.1 編集可能な結果テーブル

```
┌──────┬──────────┬──────────────────┬──────────┐
│  id  │  name    │  email           │          │
├──────┼──────────┼──────────────────┼──────────┤
│  1   │  Alice   │  alice@test.com  │  [🗑 削除]│
│  2   │  Bob     │  bob@test.com    │  [🗑 削除]│
│  3   │ *Carol*  │ *carol@new.com*  │  [🗑 削除]│  ← 編集済み（ハイライト）
├──────┼──────────┼──────────────────┼──────────┤
│ (新) │  [     ] │  [             ] │          │  ← 新規行
└──────┴──────────┴──────────────────┴──────────┘
          [変更を保存]  [変更を破棄]                ← 未保存変更がある時のみ表示
```

- PK カラムのセルは **編集不可**（グレー背景で視覚的に区別）
- 編集済みセルは **背景色をハイライト**（濃いオレンジ系）
- 新規行は最下部に表示
- 変更がある場合のみ「変更を保存」「変更を破棄」ボタンを表示

### 2.2 セル編集

- セルを **ダブルクリック** で編集モードに入る
- 編集中はセルが `<input>` に変わる
- **Enter** で確定、**Escape** でキャンセル
- **Tab** で次のセルに移動
- 空文字を入力した場合は空文字として保存（NULL にはしない）
- NULL を入力するには、セル内容を全削除してから Escape するのではなく、明示的な手段が必要
  → セル右クリック or 専用キー（Ctrl+Shift+N）で NULL にセット

### 2.3 行の追加

- テーブル最下部に空行を1行表示
- 空行にデータを入力して「変更を保存」で INSERT
- PK カラムが AUTOINCREMENT の場合、PK セルは空のまま（DB が自動採番）

### 2.4 行の削除

- 各行の右端に削除ボタン（🗑）を表示
- クリックで行を「削除予定」としてマーク（取り消し線 + 赤背景）
- 「変更を保存」で実際に DELETE 実行

### 2.5 変更の保存

「変更を保存」クリック時:

1. **確認ダイアログ** を表示
   ```
   以下の変更を実行します:
   - UPDATE 2件
   - INSERT 1件
   - DELETE 1件

   [キャンセル] [実行]
   ```
2. 確認後、生成した SQL をまとめて `POST /api/exec/batch` で送信
3. バックエンドがトランザクション内で全 SQL を実行（全成功 or 全ロールバック）
4. 成功したら結果を再取得（元の SELECT を再実行）
5. エラーが出たらロールバック済みなのでエラー表示のみ

### 2.6 変更の破棄

「変更を破棄」クリックで全ての未保存変更をリセットし、元のデータに戻す。

---

## 3. SQL 自動生成

### 3.1 UPDATE

```sql
UPDATE <tableName> SET <col1> = '<val1>', <col2> = '<val2>'
WHERE <pk1> = '<pkVal1>' AND <pk2> = '<pkVal2>'
```

- 変更されたカラムのみ SET に含める
- 値は適切にエスケープする（シングルクォート内のシングルクォートは `''` に）
- NULL は `SET col = NULL`（クォートなし）

### 3.2 INSERT

```sql
INSERT INTO <tableName> (<col1>, <col2>, ...) VALUES ('<val1>', '<val2>', ...)
```

- PK カラムの値が空の場合は INSERT 文から除外（AUTOINCREMENT 対応）
- NULL 値は `NULL`（クォートなし）

### 3.3 DELETE

```sql
DELETE FROM <tableName> WHERE <pk1> = '<pkVal1>' AND <pk2> = '<pkVal2>'
```

---

## 4. バックエンド変更

### 4.1 `POST /api/exec/batch` — バッチ実行 API

複数の SQL をトランザクション内で一括実行する。

**Request:**
```json
{
  "statements": [
    "UPDATE users SET name = 'Alice2' WHERE id = 1",
    "INSERT INTO users (name) VALUES ('Charlie')",
    "DELETE FROM users WHERE id = 2"
  ]
}
```

**Response（HTTP 200）:**
```json
{
  "totalAffectedRows": 3,
  "executionTimeMs": 12
}
```

**Response（エラー / HTTP 400）:**
```json
{
  "error": "statement 2: UNIQUE constraint failed: users.id"
}
```

- `BEGIN` → 全 SQL 実行 → `COMMIT`
- いずれかの SQL が失敗した場合 `ROLLBACK` してエラーを返す
- エラーメッセージには何番目の SQL で失敗したかを含める

### 4.2 `/api/query` レスポンス拡張

```go
// internal/query/types.go
type Response struct {
    Columns       []string `json:"columns"`
    Rows          [][]any  `json:"rows"`
    RowCount      int      `json:"rowCount"`
    ExecutionTime int64    `json:"executionTimeMs"`
    Truncated     bool     `json:"truncated,omitempty"`
    EditableTable string   `json:"editableTable,omitempty"` // 追加
    PKColumns     []string `json:"pkColumns,omitempty"`     // 追加
}
```

### 4.2 Database インターフェース拡張

`ExecBatch` メソッドを追加する:

```go
type Database interface {
    // ... 既存メソッド
    ExecBatch(statements []string) (totalAffected int64, err error)
}
```

各ドライバで `sql.Tx` を使って BEGIN → 全実行 → COMMIT / ROLLBACK を実装する。
共通ヘルパー `execBatch(db *sql.DB, statements []string)` を `scan.go` に追加。

### 4.3 単一テーブル判定 + PK 取得

`query/service.go` に判定ロジックを追加:

```go
func DetectEditableTable(sql string, db database.Database) (tableName string, pkColumns []string) {
    // 1. 正規表現で FROM <tableName> を抽出
    // 2. JOIN, UNION, サブクエリがあれば空を返す
    // 3. db.Columns(tableName) で PK カラムを取得
    // 4. PK が SELECT の columns に含まれるか確認
}
```

### 4.3 ディレクトリ変更

```
internal/query/
├── handler.go        # 既存（レスポンスに editableTable, pkColumns を追加）
├── service.go        # 既存（DetectEditableTable 追加）
├── types.go          # 既存（Response に 2 フィールド追加）
├── detect.go         # 新規: テーブル名判定ロジック
└── detect_test.go    # 新規: 判定ロジックのテスト
```

---

## 5. フロントエンド変更

### 5.1 新規・変更ファイル

```
frontend/src/
├── features/
│   ├── results-panel/
│   │   ├── ResultsPanel.vue       # 変更: 編集可能テーブルの分岐
│   │   ├── EditableTable.vue      # 新規: 編集可能テーブルコンポーネント
│   │   ├── EditableCell.vue       # 新規: セル編集コンポーネント
│   │   ├── ConfirmDialog.vue      # 新規: 変更確認ダイアログ
│   │   ├── use-table-edit.ts      # 新規: 編集状態管理 composable
│   │   ├── sql-generator.ts       # 新規: UPDATE/INSERT/DELETE 文生成
│   │   └── types.ts               # 変更: 編集関連の型追加
├── shared/
│   └── types.ts                   # 変更: QueryResult に editableTable, pkColumns 追加
```

### 5.2 shared/types.ts 変更

```typescript
export type QueryResult = {
  columns: string[];
  rows: (string | null)[][];
  rowCount: number;
  executionTimeMs: number;
  truncated?: boolean;
  editableTable?: string;   // 追加
  pkColumns?: string[];     // 追加
};
```

### 5.3 use-table-edit.ts

編集状態を管理する composable:

```typescript
type CellEdit = {
  rowIndex: number;
  colIndex: number;
  originalValue: string | null;
  newValue: string | null;
};

type NewRow = {
  values: (string | null)[];
};

type DeletedRow = {
  rowIndex: number;
};
```

- `edits: Ref<CellEdit[]>` — セル編集のリスト
- `newRows: Ref<NewRow[]>` — 追加行のリスト
- `deletedRows: Ref<Set<number>>` — 削除対象行のインデックス
- `hasChanges: ComputedRef<boolean>` — 未保存変更の有無
- `updateCell(rowIndex, colIndex, value)` — セル値を更新
- `addRow()` — 空行を追加
- `deleteRow(rowIndex)` — 行を削除マーク
- `discardAll()` — 全変更を破棄
- `generateSQL()` — 全変更の SQL を生成

---

## 6. テスト

### 6.1 バックエンド

| テストファイル | テスト内容 |
|--------------|-----------|
| `query/detect_test.go` | 単一テーブル判定（正常系・JOIN・サブクエリ・UNION） |
| `query/handler_test.go` | editableTable / pkColumns がレスポンスに含まれる |

### 6.2 フロントエンド E2E

| テスト | 内容 |
|--------|------|
| 編集可能テーブル表示 | SELECT * FROM でセルがダブルクリック編集可能 |
| セル編集 + 保存 | 値を変更 → 保存 → 再取得で反映される |
| 行追加 + 保存 | 新規行に入力 → 保存 → SELECT で追加される |
| 行削除 + 保存 | 削除ボタン → 保存 → SELECT から消える |
| 変更破棄 | 編集後「変更を破棄」→ 元の値に戻る |
| PK なしテーブル | PK なしテーブルの SELECT → 読み取り専用 |

---

## 7. Issue 一覧（Phase 5）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 32 | 単一テーブル判定 + PK 検出 | query/detect.go | なし |
| 33 | /api/query レスポンス拡張 | query/service.go, types.go, handler.go | #32 |
| 34 | SQL 自動生成（UPDATE/INSERT/DELETE） | results-panel/sql-generator.ts | なし |
| 35 | 編集可能テーブル UI | EditableTable.vue, EditableCell.vue | #33, #34 |
| 36 | 編集状態管理 composable | use-table-edit.ts | #34 |
| 37 | 確認ダイアログ + 保存/破棄 | ConfirmDialog.vue, ResultsPanel.vue | #35, #36 |
| 38 | テスト追加（Phase 5 機能） | detect_test.go, e2e/ | #37 |

---

## 8. Phase 5 完了条件

### サーバー
- [ ] `/api/query` が editableTable と pkColumns をレスポンスに含む
- [ ] 単一テーブル SELECT のみ editableTable が返される
- [ ] JOIN / UNION / サブクエリの場合は editableTable が空
- [ ] PK がないテーブルは pkColumns が空
- [ ] `go test ./internal/...` で全テストがパスする

### フロントエンド
- [ ] 編集可能な結果テーブルでセルをダブルクリック編集できる
- [ ] PK カラムは編集不可
- [ ] 変更されたセルがハイライト表示される
- [ ] 新規行を追加して INSERT できる
- [ ] 行を削除して DELETE できる
- [ ] 「変更を保存」で確認ダイアログ → SQL 実行 → 結果再取得
- [ ] 「変更を破棄」で全変更がリセットされる
- [ ] PK なしテーブルでは読み取り専用のまま
- [ ] `pnpm build` が成功する
- [ ] `pnpm test:e2e` で全テストがパスする
