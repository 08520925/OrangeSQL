# フロントエンド（frontend/）

## アーキテクチャ

Vue 3 + TypeScript + Vite。機能別ディレクトリ（features/）に Vue SFC + composable + 型をコロケーション。

```
src/
├── App.vue                 ルートコンポーネント（全体レイアウト + 状態統合）
├── main.ts                 エントリポイント
├── shared/
│   ├── api.ts              API クライアント（全エンドポイント）
│   └── types.ts            QueryResult, ExecResult, BatchExecResult 等
├── features/
│   ├── sql-editor/         CodeMirror 6 エディタ
│   │   ├── SqlEditor.vue
│   │   ├── use-sql-editor.ts     EditorView 管理
│   │   ├── sql-completions.ts    オートコンプリート（テーブル・カラム・キーワード）
│   │   └── types.ts
│   ├── results-panel/      結果テーブル + 編集 + エクスポート
│   │   ├── ResultsPanel.vue      読み取り専用 / 編集可能の分岐
│   │   ├── EditableTable.vue     インライン編集 UI
│   │   ├── EditableCell.vue      セル編集コンポーネント
│   │   ├── ExportToolbar.vue     CSV/JSON/Markdown エクスポート
│   │   ├── use-table-edit.ts     編集状態管理
│   │   ├── sql-generator.ts      UPDATE/INSERT/DELETE 文生成
│   │   ├── result-exporter.ts    CSV/JSON/Markdown 変換 + ダウンロード/コピー
│   │   └── types.ts
│   ├── schema-sidebar/     サイドバー（テーブル一覧 + カラム展開）
│   ├── connection/         接続プロファイル管理
│   │   ├── ConnectionDropdown.vue
│   │   ├── ConnectionDialog.vue   ドライバ別フォーム切り替え
│   │   ├── ConnectionManager.vue
│   │   ├── use-connection.ts
│   │   └── types.ts              DriverType, DRIVERS, DEFAULT_PORTS
│   ├── tab-bar/            タブ管理
│   └── resize-handle/      パネルリサイズ
└── e2e/                    Playwright E2E テスト（phase1〜6）
```

## ルール

- TypeScript 最厳格設定: `strict` + `noUncheckedIndexedAccess` + `exactOptionalPropertyTypes`
- `any` 禁止、`as` キャスト禁止 → `unknown` + 型ガードを使う
- Vue SFC からロジックを composable（`use-*.ts`）に抽出
- CodeMirror 6 は Vue ラッパー禁止、素で `new EditorView()` を使う
- API クライアントは `shared/api.ts` に集約（直接 fetch しない）
- コンポーネント間の通信は props + emit（グローバル state は使わない）
- CSS は scoped、カラーテーマはダークベース（#1e1e1e 系）

## API プロキシ（開発時）

`vite.config.ts` で `/api` を `http://localhost:5522` にプロキシ。
プロダクションでは Go サーバーが API + 静的ファイル両方を配信するためプロキシ不要。

## E2E テスト

- Playwright 使用
- バックエンド（:5522）+ フロントエンド dev server（:5173）を起動してから実行
- テストデータは API 経由で直接 INSERT してセットアップ
- `pnpm test:e2e` で全テスト実行
