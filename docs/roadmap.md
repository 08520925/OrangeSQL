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

### Phase 5: テーブルデータ編集
- 結果テーブルのセルをダブルクリックでインライン編集
- UPDATE / INSERT / DELETE 文をフロントで自動生成
- `POST /api/exec/batch` でトランザクション一括実行（全成功 or 全ロールバック）
- PK 必須（PK がないテーブルは読み取り専用）
- 単一テーブル SELECT のみ編集可能（JOIN / UNION / GROUP BY 等は読み取り専用）
- 変更確認ダイアログ + 変更破棄機能
- 仕様書: `docs/spec/phase5-table-edit.md`

### Phase 6: オートコンプリート + エクスポート
- SQL オートコンプリート（テーブル名・カラム名・SQL キーワード）
- `GET /api/schema/completions` で全テーブル+カラムを一括取得
- CodeMirror 6 `autocompletion` 拡張で補完ポップアップ
- `table.` でカラム名補完、Ctrl+Space で手動トリガー
- 結果エクスポート: CSV / JSON ダウンロード
- クリップボードコピー: CSV / Markdown テーブル形式
- 仕様書: `docs/spec/phase6-autocomplete-export.md`

### Phase 7: シングルバイナリ配布
- `go:embed` でフロントエンド（frontend/dist）を exe に埋め込み
- exe 単体で起動 → API + 静的ファイル配信（dev サーバー不要）
- 起動時にデフォルトブラウザ自動オープン（`--no-browser` で抑制）
- SPA フォールバック対応
- `build.sh` でフロントビルド → Go ビルドを一発実行
- 仕様書: `docs/spec/phase7-single-binary.md`

---

## 検討中（Phase 未定）

- クエリ履歴
- タブの永続化（ブラウザリロードで復元）
- ライト/ダークテーマ切り替え
- サイドバー幅リサイズ
- タブタイトルの編集
