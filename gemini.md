# OrangeSQL – Gemini CLI Project Guide

## ルール
回答は必ず日本語で回答すること。

## 1. プロジェクト概要
OrangeSQL は Wails v2 を使ったデスクトップ DB クライアント（A5:SQL 風）です。
フロントは Vue 3 + TypeScript + Vite、SQL エディタは Monaco Editor を採用します。
最終成果物は Windows 向けの単一 exe 配布を前提とします（オフライン動作）。

## 2. 重要方針（必ず守る）
- **最終的に exe 配布前提**：フロントは `frontend/dist` を生成し、Go 側で embed して配布する。
- **CGO は基本 0**：将来的に PostgreSQL / MySQL / SQL Server を対象にするため、基本は pure Go ドライバ方針。
  - SQLite も cgo 不要ドライバ（例：modernc）を優先検討。
- **Monaco は Vue ラッパーを使わず「素で初期化」**：依存衝突や bundle 問題を回避する。
- **CSS でルートに `text-align:center` を置かない**：Monaco/Canvas系が崩れるため禁止。
- Wails の UI は WebView2 上で動くため、ブラウザと違う制約がある（DevTools、右クリックなど）。

## 3. 技術スタック
- Backend: Go (Wails v2)
- Frontend: Vue 3 + TypeScript + Vite
- SQL editor: Monaco Editor（`monaco-editor` を直接利用）
- DB:
  - SQLite（pure Go を優先）
  - PostgreSQL / MySQL / SQL Server（将来対応）

## 4. ディレクトリ構成（現在の想定）
- `/main.go`：Wails 起動・assets embed・Bind 設定
- `/app.go`：Wails に Bind される窓口（フロントから呼べる API）
- `/frontend/`：Vue + Vite
  - `/frontend/src/`：UI
  - `/frontend/dist/`：ビルド成果物（Wails が埋め込む）
- `/internal/`（将来）：アプリ層 / サービス層 / DB 層などを分離する予定

## 5. よく使うコマンド
### 開発起動
- `wails dev`

### ビルド（exe）
- `wails build`

### フロント単体
- `cd frontend`
- `npm run dev`
- `npm run build`

## 6. Monaco Editor 実装ルール（重要）
- Vue ラッパー禁止：`@guolao/vue-monaco-editor` を導入しない（peer conflict の原因になる）
- Monaco は `onMounted` で `monaco.editor.create()` を実行する
- 親要素は必ず高さを持つ（`height:100vh` か flex + `min-height:0`）
- 初期化直後に `setTimeout(() => editor.layout(), 0)` を呼ぶ（WebView2 対策）
- ルート CSS に `text-align:center` を入れない

## 7. DB 方針（将来拡張を見据える）
- `database/sql` を共通インターフェースとして使う
- DBごとに driver と接続文字列を切り替える設計にする
- SQL 実行 API は「Query（SELECT）」と「Exec（更新系）」を分ける
- 結果は GUI 表示向けに `columns: string[]` と `rows: string[][]` を基本フォーマットとする

## 8. 作業依頼の仕方（Gemini への指示規約）
あなた（Gemini）は以下を優先して提案・実装してください：
1) **exe配布で壊れない**（外部 CDN を避ける / dist に同梱）
2) **依存関係の衝突を避ける**（特に Monaco 周辺は慎重）
3) **Windows + WebView2 前提**の挙動（DevTools/右クリック/IME/フォーカス）
4) 変更は最小・段階的（まず動く最小 → 機能追加）

## 9. いまの優先タスク（ロードマップ）
- [ ] SQL エディタ（Monaco）を安定表示
- [ ] Ctrl+Enter で SQL 実行
- [ ] 結果グリッド表示（列ヘッダ固定）
- [ ] テーブル一覧ペイン（SQLite から）
- [ ] 接続プロファイル管理（将来的に PG/MySQL/MSSQL）
- [ ] exe build 確認・サイズ最適化

## 10. コードを書くときの約束
- TypeScript は strict 前提
- UI は最初は inline style でもよいが、後でコンポーネント化しやすい構造にする
- 例外はユーザーに分かるメッセージで返す（Go→JSの error を整形）
- 迷ったら「A5:SQL っぽい UX（軽快・即実行・履歴）」を優先する

## Go アーキテクチャ方針

このプロジェクトは「軽量レイヤードアーキテクチャ」を採用する。

- app.go は Wails Bind 専用（Controller 相当）
- ユースケースは internal/app に配置する
- DB 抽象は domain.Database interface として定義する
- DB ごとの差分は infra 配下で吸収する
- app.go から infra を直接触らない
- database/sql を直接使ってよいのは infra 層のみ

目的：
- DB 種別追加時の影響範囲を最小化する
- UI（Vue/Wails）と DB 実装を疎結合に保つ

## Frontend Architecture (Feature-Sliced Design)

採用: Feature-Sliced Design (FSD)

- Layers: app / pages / widgets / features / entities / shared（processesは使わない）
- Segment: ui / model / api / lib / config を基本とする
- Import Rule: 上位レイヤー → 下位レイヤーのみ。feature→feature の直接importは禁止。
- Public API: 各sliceは必ず `index.ts` を持ち、外部は `index.ts` 経由でのみ import する。

目的:
- 依存方向を固定して保守性を上げる
- DB種別追加やUI増加時も影響範囲を局所化する
