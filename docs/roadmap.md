# OrangeSQL ロードマップ

## 完了済み

### Phase 1: 基盤
- Go Webサーバー + Vue 3 フロントエンド
- SQLite 接続、SQL 実行、結果テーブル表示
- REST API（query / exec / schema / health）
- CodeMirror 6 エディタ（Ctrl+Enter 実行、カーソル位置のSQL文実行）
- Playwright E2E テスト
- 仕様書: `docs/spec/phase1-foundation.md`

### Phase 2: 日常使い
- 複数タブ（追加・切り替え・閉じ）
- カラム情報表示（テーブル展開で型・PK・NOT NULL）
- パネルリサイズ（エディタ/結果パネルのドラッグ分割）
- 仕様書: `docs/spec/phase2-daily-use.md`

### Phase 3: 接続管理
- 接続プロファイル管理（作成・編集・削除）
- プロファイル切り替え（ドロップダウンメニュー）
- profiles.json で永続化（~/.orangesql/）
- 仕様書: `docs/spec/phase3-connection.md`

### Phase 4: マルチDB対応
- PostgreSQL ドライバ追加（`github.com/jackc/pgx/v5`）
- MySQL ドライバ追加（`github.com/go-sql-driver/mysql`）
- SQL Server ドライバ追加（`github.com/microsoft/go-mssqldb`）
- Database ファクトリ（`database.New(driver, params)`）で接続を抽象化
- scanRows / execStatement 共通化
- スキーマ取得（テーブル一覧・カラム情報）を DB 毎に対応
- 接続ダイアログでドライバ切り替え → フォームフィールド動的変更
- Profile に TCP 系接続パラメータ（host, port, user, password, dbName, sslMode）を追加
- ドライバ別バリデーション
- 仕様書: `docs/spec/phase4-multi-db.md`

---

## 今後

### Phase 5: テーブルデータ編集
結果グリッド上で直接データを編集できるようにする。

- 結果テーブルのセルをダブルクリックでインライン編集
- 編集内容から UPDATE 文を自動生成して実行
- 行の追加（INSERT 生成）
- 行の削除（DELETE 生成）
- PK が必要（PK がないテーブルは編集不可）
- 編集前の確認ダイアログ or Undo 機能

---

## 検討中（Phase 未定）

- クエリ履歴
- タブの永続化（ブラウザリロードで復元）
- ライト/ダークテーマ切り替え
- サイドバー幅リサイズ
- タブタイトルの編集
- Go embed によるシングルバイナリ配布
- クエリ結果の CSV / JSON エクスポート
- SQL オートコンプリート（テーブル名・カラム名の補完）
