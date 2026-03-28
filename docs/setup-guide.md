# OrangeSQL セットアップ・運用ガイド

このドキュメントは、開発を開始・継続するために必要な手順をまとめたものです。

---

## 1. 初回セットアップ（未完了の作業）

### 1.1 GitHub に ANTHROPIC_API_KEY を設定する

Claude Code GitHub Action を動かすために必要。

1. https://github.com/08520925/OrangeSQL/settings/secrets/actions にアクセス
2. 「New repository secret」をクリック
3. Name: `ANTHROPIC_API_KEY`
4. Value: Anthropic の API キーを貼り付け
5. 「Add secret」をクリック

### 1.2 Claude GitHub App をインストールする

Issue への `@claude` メンションで自動実装を有効にする。

1. ターミナルで Claude Code を起動
2. `/install-github-app` を実行
3. 画面の指示に従ってリポジトリにアプリをインストール

### 1.3 現在の変更をコミット・プッシュする

仕様書・CLAUDE.md・GitHub Action ワークフローをリポジトリに反映する。

```bash
cd C:/work/OrangeSQL
git add docs/ CLAUDE.md .github/
git commit -m "Phase 1 仕様書・CLAUDE.md・GitHub Action 追加"
git push origin develop
```

### 1.4 pnpm をインストールする（未インストールの場合）

```bash
npm install -g pnpm
```

### 1.5 Go 1.26.1 を確認する

```bash
go version
# go1.26.1 でなければ https://go.dev/dl/ からインストール
```

---

## 2. 開発の進め方

### 方法A: スマホから Issue 駆動（推奨）

1. GitHub の Issue 一覧を開く: https://github.com/08520925/OrangeSQL/issues
2. 着手したい Issue を開く（依存関係に注意、#1 が最初）
3. Issue にコメント: `@claude この Issue を実装して`
4. Claude が自動でブランチ作成 → 実装 → PR 作成
5. PR を確認してマージ

### 方法B: ローカルで Claude Code を使う

```bash
cd C:/work/OrangeSQL
claude
# プロンプト例: "Issue #1 を実装して"
```

### Issue の実装順序（依存関係）

```
#1 プロジェクト基盤（最初にやる）
 ├── #2 POST /api/query
 ├── #3 POST /api/exec
 ├── #4 GET /api/schema/tables
 ├── #5 GET /api/health
 └── #6 フロントエンド基盤
      ├── #7 SQLエディタ
      ├── #8 結果パネル（#2 も必要）
      └── #9 スキーマサイドバー（#4 も必要）
           └── #10 統合（#7, #8, #9 全て必要）
```

---

## 3. 開発時のコマンド

### バックエンド

```bash
cd OrangeSQL

# サーバー起動（デフォルト: data.db）
go run .

# 別の DB ファイルを指定
go run . -db /path/to/mydb.sqlite

# テスト実行
go test ./internal/...
```

### フロントエンド

```bash
cd OrangeSQL/frontend

# 依存インストール
pnpm install

# 開発サーバー起動（:5173）
pnpm dev

# ビルド
pnpm build

# パッケージ追加
pnpm add <package>
```

### 開発時の起動手順（バックエンド + フロントエンド）

ターミナルを2つ開く:

```bash
# ターミナル1: バックエンド
cd OrangeSQL && go run .

# ターミナル2: フロントエンド
cd OrangeSQL/frontend && pnpm dev
```

ブラウザで http://localhost:5173 を開く（Vite が API を :5522 にプロキシする）。

---

## 4. 関連リンク

| リンク | 内容 |
|--------|------|
| [Issue 一覧](https://github.com/08520925/OrangeSQL/issues) | Phase 1 の全タスク |
| [仕様書](docs/spec/phase1-foundation.md) | Phase 1 の詳細仕様 |
| [CLAUDE.md](CLAUDE.md) | AI 向けプロジェクトルール |
