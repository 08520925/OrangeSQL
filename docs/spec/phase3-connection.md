# Phase 3: 接続管理 — 仕様書

## 概要

複数の SQLite データベースを接続プロファイルとして管理し、UI から切り替えられるようにする。
Phase 4 でマルチ DB 対応する際の基盤を構築する。

---

## 1. 接続プロファイル

### 1.1 プロファイルの定義

```typescript
type ConnectionProfile = {
  id: string;
  name: string;          // 表示名（例: "開発DB", "本番DB"）
  driver: string;        // "sqlite"（Phase 4 で "postgres", "mysql", "sqlserver" を追加）
  path: string;          // SQLite ファイルパス
  createdAt: string;     // ISO 8601
};
```

### 1.2 プロファイルの保存先

- アプリのデータディレクトリに JSON ファイルとして保存
- パス: `~/.orangesql/profiles.json`（Windows: `%USERPROFILE%\.orangesql\profiles.json`）
- DB 内ではなくファイルに保存する（DB を切り替えても設定が消えない）

### 1.3 profiles.json の形式

```json
{
  "profiles": [
    {
      "id": "1",
      "name": "開発DB",
      "driver": "sqlite",
      "path": "C:/work/project/dev.db",
      "createdAt": "2026-03-28T12:00:00Z"
    }
  ],
  "lastUsedId": "1"
}
```

---

## 2. UI

### 2.1 ヘッダーバーの変更

```
┌────────────────────────────────────────────────────┐
│  OrangeSQL  [開発DB ▼]                    ▶ 実行   │
└────────────────────────────────────────────────────┘
```

- 現在の DB 名表示部分をドロップダウンに変更
- クリックでプロファイル一覧が表示される
- ドロップダウン内に「+ 新規接続」ボタン

### 2.2 ドロップダウンメニュー

```
┌──────────────────┐
│ ● 開発DB         │  ← 現在の接続（チェックマーク）
│   本番DB         │
│   テスト用DB     │
├──────────────────┤
│ + 新規接続       │
│ ⚙ 接続管理      │
└──────────────────┘
```

- プロファイルをクリックで接続切り替え
- 「+ 新規接続」で接続ダイアログを開く
- 「⚙ 接続管理」で接続一覧管理ダイアログを開く

### 2.3 接続ダイアログ（新規 / 編集）

```
┌─── 新規接続 ──────────────────┐
│                               │
│  接続名:  [開発DB           ] │
│                               │
│  ドライバ: [SQLite ▼]        │
│                               │
│  ファイルパス:                 │
│  [C:/work/dev.db     ] [参照] │
│                               │
│      [キャンセル]  [接続]     │
└───────────────────────────────┘
```

- 接続名: 任意のテキスト（必須）
- ドライバ: Phase 3 では SQLite 固定（ドロップダウンは表示するが選択肢は1つ）
- ファイルパス: テキスト入力（必須）。「参照」ボタンは Phase 3 ではなし（テキスト入力のみ）
- 「接続」ボタンで保存 + 即座に接続切り替え

### 2.4 接続管理ダイアログ

```
┌─── 接続管理 ──────────────────┐
│                               │
│  開発DB     sqlite  [編集][削除]│
│  本番DB     sqlite  [編集][削除]│
│  テストDB   sqlite  [編集][削除]│
│                               │
│               [閉じる]        │
└───────────────────────────────┘
```

- 各プロファイルの編集・削除
- 現在接続中のプロファイルは削除不可
- 「編集」で接続ダイアログを開く（既存値がセットされた状態）

---

## 3. API

### 3.1 プロファイル管理 API

#### `GET /api/profiles`

プロファイル一覧を取得する。

**Response（HTTP 200）:**
```json
{
  "profiles": [
    { "id": "1", "name": "開発DB", "driver": "sqlite", "path": "C:/work/dev.db", "createdAt": "2026-03-28T12:00:00Z" }
  ],
  "activeId": "1"
}
```

#### `POST /api/profiles`

新しいプロファイルを作成する。

**Request:**
```json
{
  "name": "テストDB",
  "driver": "sqlite",
  "path": "C:/work/test.db"
}
```

**Response（HTTP 201）:**
```json
{
  "id": "2",
  "name": "テストDB",
  "driver": "sqlite",
  "path": "C:/work/test.db",
  "createdAt": "2026-03-28T13:00:00Z"
}
```

#### `PUT /api/profiles/{id}`

プロファイルを更新する。

**Request:**
```json
{
  "name": "テストDB（更新）",
  "path": "C:/work/test2.db"
}
```

**Response（HTTP 200）:** 更新後のプロファイル。

#### `DELETE /api/profiles/{id}`

プロファイルを削除する。現在接続中のプロファイルは削除不可（HTTP 400）。

**Response（HTTP 200）:**
```json
{ "deleted": true }
```

#### `POST /api/profiles/{id}/connect`

指定プロファイルに接続を切り替える。

**Response（HTTP 200）:**
```json
{ "connected": true, "database": "C:/work/test.db" }
```

**Response（エラー / HTTP 400）:**
```json
{ "error": "failed to open database: no such file" }
```

---

## 4. バックエンド変更

### 4.1 ディレクトリ追加

```
internal/
├── profile/                   # 接続プロファイル管理
│   ├── handler.go             # REST API ハンドラ
│   ├── handler_test.go
│   ├── service.go             # プロファイル CRUD + 接続切り替え
│   ├── storage.go             # profiles.json の読み書き
│   └── types.go               # Profile 型定義
```

### 4.2 接続切り替えの仕組み

- `main.go` で現在の DB 接続を `*database.SQLite` として保持
- `POST /api/profiles/{id}/connect` 時に:
  1. 新しい DB に接続テスト
  2. 成功したら旧 DB を Close
  3. 全ハンドラの DB 参照を新しい DB に差し替え
- DB 参照は `sync.RWMutex` で保護する

```go
// internal/profile/service.go
type ConnectionManager struct {
    mu      sync.RWMutex
    current database.Database
    storage *Storage
}

func (m *ConnectionManager) DB() database.Database {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.current
}

func (m *ConnectionManager) SwitchTo(profile Profile) error {
    // 新DB接続 → 旧DB切断 → 差し替え
}
```

各ハンドラは `ConnectionManager.DB()` 経由で DB にアクセスする。

### 4.3 起動時の動作

1. `~/.orangesql/profiles.json` を読み込む
2. `lastUsedId` のプロファイルに接続
3. profiles.json が存在しない場合:
   - `-db` フラグのパス（デフォルト `./data.db`）で初期プロファイルを自動作成
   - profiles.json を生成

---

## 5. フロントエンド変更

### 5.1 新規ファイル

```
frontend/src/
├── features/
│   ├── connection/            # 接続管理
│   │   ├── ConnectionDropdown.vue    # ヘッダーのドロップダウン
│   │   ├── ConnectionDialog.vue      # 新規/編集ダイアログ
│   │   ├── ConnectionManager.vue     # 接続管理ダイアログ
│   │   ├── use-connection.ts         # 接続管理 composable
│   │   └── types.ts                  # ConnectionProfile 型
```

### 5.2 shared/api.ts 追加

```typescript
export function fetchProfiles(): Promise<ProfilesResponse>;
export function createProfile(data: CreateProfileRequest): Promise<Profile>;
export function updateProfile(id: string, data: UpdateProfileRequest): Promise<Profile>;
export function deleteProfile(id: string): Promise<void>;
export function connectProfile(id: string): Promise<{ connected: boolean; database: string }>;
```

### 5.3 接続切り替え後の動作

- サイドバーのテーブル一覧を再取得
- 全タブの結果をクリア（idle に戻す）
- ヘッダーの DB 名を更新

---

## 6. テスト

### 6.1 バックエンド

| テストファイル | テスト内容 |
|--------------|-----------|
| `profile/handler_test.go` | プロファイル CRUD、接続切り替え、削除制約 |
| `profile/storage_test.go` | profiles.json の読み書き |

### 6.2 E2E

| テスト | 内容 |
|--------|------|
| プロファイル作成 | 新規接続ダイアログで作成 → ドロップダウンに追加される |
| 接続切り替え | ドロップダウンから別のプロファイルに切り替え → サイドバー更新 |
| プロファイル削除 | 接続管理から削除 → ドロップダウンから消える |

---

## 7. Issue 一覧（Phase 3）

| # | タイトル | スコープ | 依存 |
|---|---------|---------|------|
| 17 | プロファイル保存・読込（profiles.json） | profile/storage.go | なし |
| 18 | ConnectionManager（接続切り替え機構） | profile/service.go | #17 |
| 19 | プロファイル REST API | profile/handler.go | #18 |
| 20 | main.go をプロファイル対応に改修 | main.go | #18 |
| 21 | フロント: 接続ドロップダウン + ダイアログ | features/connection/ | #19 |
| 22 | E2E テスト追加（Phase 3 機能） | e2e/ | #21 |

---

## 8. Phase 3 完了条件

### サーバー
- [ ] `~/.orangesql/profiles.json` にプロファイルが保存される
- [ ] プロファイルの作成・編集・削除が API で動作する
- [ ] `POST /api/profiles/{id}/connect` で接続が切り替わる
- [ ] 現在接続中のプロファイルは削除できない
- [ ] 起動時に前回のプロファイルに自動接続する
- [ ] `go test ./internal/...` で全テストがパスする

### フロントエンド
- [ ] ヘッダーにプロファイルドロップダウンが表示される
- [ ] ドロップダウンから接続を切り替えられる
- [ ] 新規接続ダイアログでプロファイルを作成できる
- [ ] 接続管理ダイアログでプロファイルを編集・削除できる
- [ ] 接続切り替え後にサイドバーが更新される
- [ ] `pnpm build` が成功する
- [ ] `pnpm test:e2e` で全テストがパスする
