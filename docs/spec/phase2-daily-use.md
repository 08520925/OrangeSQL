# Phase 2: 日常使い — 仕様書

## 概要

OrangeSQL を日常のSQL作業で使えるレベルに引き上げる。
複数クエリの同時作業、テーブル構造の確認、画面カスタマイズを実現する。

---

## 1. 複数タブ

### 1.1 タブの動作

```
┌──────┬──────┬──────┬───┐
│ Tab1 │ Tab2 │ Tab3 │ + │
├──────┴──────┴──────┴───┤
│  SQLエディタ（アクティブタブの内容）│
```

- ヘッダーバーの直下にタブバーを表示する
- 各タブは独立した SQL エディタと結果パネルを持つ
- タブ切り替え時にエディタの内容と結果を保持する
- 「+」ボタンで新しいタブを追加（上限なし）
- タブに「×」ボタンで閉じる（最後の1タブは閉じられない）

### 1.2 タブの状態

各タブが保持する状態:

```typescript
type Tab = {
  id: string;
  title: string;         // デフォルト: "Query 1", "Query 2", ...
  sql: string;           // エディタの内容
  result: ResultState;   // 実行結果（idle/loading/query/exec/error）
};
```

### 1.3 タブの表示

- アクティブタブはハイライト（背景色を明るく）
- タブ幅は内容に応じて可変（最小80px、最大200px）
- タブが多い場合は水平スクロール
- タブタイトルはダブルクリックで編集可能（Phase 2 では未実装、Phase 3 以降）

### 1.4 初期状態

- アプリ起動時は1タブ（"Query 1"）を表示
- タブの状態はメモリ上のみ保持（ブラウザリロードでリセット）

---

## 2. カラム情報表示

### 2.1 テーブル展開

サイドバーのテーブル名をクリックすると、従来通り SELECT 文をエディタに挿入する。
テーブル名の左にある **▶ アイコン** をクリックすると、カラム一覧を展開表示する。

```
▼ ▦ users
   🔑 id       INTEGER  PK
      name     TEXT     NOT NULL
      email    TEXT
▶ ▦ orders
```

### 2.2 表示する情報

| 項目 | 表示 |
|------|------|
| カラム名 | そのまま表示 |
| データ型 | `INTEGER`, `TEXT`, `REAL`, `BLOB` 等 |
| PK | 🔑 アイコン |
| NOT NULL | ラベル表示 |

### 2.3 API

#### `GET /api/schema/columns?table={tableName}`

指定テーブルのカラム情報を取得する。

**Response（HTTP 200）:**
```json
{
  "columns": [
    { "name": "id", "type": "INTEGER", "pk": true, "notNull": true },
    { "name": "name", "type": "TEXT", "pk": false, "notNull": true },
    { "name": "email", "type": "TEXT", "pk": false, "notNull": false }
  ]
}
```

**Response（エラー / HTTP 400）:**
```json
{
  "error": "table not found: nonexistent"
}
```

### 2.4 SQLite でのカラム情報取得

```sql
PRAGMA table_info({tableName})
```

返り値: `cid`, `name`, `type`, `notnull`, `dflt_value`, `pk`

---

## 3. パネルリサイズ

### 3.1 動作

- エディタと結果パネルの境界にドラッグハンドルを表示する
- ハンドルをドラッグして上下の分割比率を変更できる
- 最小高さ: エディタ 100px、結果パネル 100px

### 3.2 ドラッグハンドル

- 高さ 4px のバー（ホバー時に色が変わる）
- カーソルが `row-resize` に変わる
- ドラッグ中は半透明のオーバーレイを表示して操作感を向上

### 3.3 分割比率の保持

- ドラッグで変更した比率はメモリ上で保持
- ブラウザリロードでデフォルト（50:50）にリセット

---

## 4. バックエンド変更

### 4.1 新規エンドポイント

| エンドポイント | メソッド | 内容 |
|--------------|---------|------|
| `/api/schema/columns` | GET | カラム情報取得 |

### 4.2 Database interface 拡張

```go
type Database interface {
    // ... 既存メソッド ...

    // Columns は指定テーブルのカラム情報を返す。
    Columns(tableName string) ([]ColumnInfo, error)
}

type ColumnInfo struct {
    Name    string
    Type    string
    PK      bool
    NotNull bool
}
```

### 4.3 ディレクトリ変更

```
internal/
├── schema/                   # 既存 + columns 追加
│   ├── handler.go            # columns ハンドラ追加
│   ├── handler_test.go       # columns テスト追加
│   ├── service.go            # columns 取得ロジック追加
│   └── types.go              # ColumnInfo 型追加
├── database/
│   ├── interface.go          # Columns メソッド + ColumnInfo 型追加
│   ├── sqlite.go             # Columns 実装追加
│   └── sqlite_test.go        # Columns テスト追加
```

---

## 5. フロントエンド変更

### 5.1 新規ファイル

```
frontend/src/
├── features/
│   ├── tab-bar/               # タブ管理
│   │   ├── TabBar.vue
│   │   ├── use-tabs.ts
│   │   └── types.ts
│   ├── resize-handle/         # パネルリサイズ
│   │   ├── ResizeHandle.vue
│   │   └── use-resize.ts
│   └── schema-sidebar/        # 既存拡張
│       ├── SchemaSidebar.vue   # カラム展開を追加
│       ├── use-schema.ts       # columns 取得を追加
│       └── types.ts            # ColumnInfo 型追加
```

### 5.2 App.vue の変更

```
┌─────────────────────────────────────────────────┐
│  ヘッダーバー                                     │
├─────────────────────────────────────────────────┤
│  タブバー  [Query 1] [Query 2] [+]               │
├────────────┬────────────────────────────────────┤
│  スキーマ   │  SQLエディタ                         │
│  サイドバー  │                                     │
│            │──── ドラッグハンドル ────             │
│  テーブル   │  結果パネル                           │
│  (展開可)  │                                     │
│            │                                     │
│            │  ステータスバー                       │
└────────────┴────────────────────────────────────┘
```

### 5.3 shared/api.ts 追加

```typescript
/** カラム情報を取得する */
export function fetchColumns(table: string): Promise<{ columns: ColumnInfo[] }> {
  return request(`/schema/columns?table=${encodeURIComponent(table)}`);
}
```

---

## 6. テスト

### 6.1 バックエンド（追加分）

| テストファイル | テスト内容 |
|--------------|-----------|
| `schema/handler_test.go` | カラム情報取得、存在しないテーブル |
| `database/sqlite_test.go` | Columns メソッドのテスト |

### 6.2 E2E（追加分）

| テスト | 内容 |
|--------|------|
| タブ追加・切り替え | 新タブ作成、タブ間のSQL独立性 |
| タブ閉じ | タブを閉じる、最後の1タブは閉じられない |
| カラム情報 | テーブル展開でカラム一覧が表示される |
| パネルリサイズ | ドラッグで分割比率が変わる |

---

## 7. Issue 一覧（Phase 2）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 11 | Database interface に Columns メソッド追加 | database/ | なし |
| 12 | GET /api/schema/columns エンドポイント実装 | schema/ | #11 |
| 13 | フロント: 複数タブ機能 | features/tab-bar/ | なし |
| 14 | フロント: カラム情報表示（サイドバー拡張） | features/schema-sidebar/ | #12 |
| 15 | フロント: パネルリサイズ | features/resize-handle/ | なし |
| 16 | E2E テスト追加（Phase 2 機能） | e2e/ | #13〜#15 |

---

## 8. Phase 2 完了条件

### サーバー
- [ ] `GET /api/schema/columns?table=xxx` でカラム情報が返る
- [ ] `go test ./internal/...` で全テストがパスする

### フロントエンド
- [ ] タブバーが表示され、タブの追加・切り替え・閉じができる
- [ ] 各タブが独立したSQL・結果を保持する
- [ ] サイドバーでテーブルを展開するとカラム情報が表示される
- [ ] エディタと結果パネルの境界をドラッグでリサイズできる
- [ ] `pnpm build` が成功する
- [ ] `pnpm test:e2e` で全テストがパスする
