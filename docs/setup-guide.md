# OrangeSQL セットアップ・運用ガイド

このドキュメントは、開発を開始・継続するために必要な手順をまとめたものです。

---

## 1. 初回セットアップ

### 1.1 pnpm をインストールする（未インストールの場合）

```bash
npm install -g pnpm
```

### 1.2 Go 1.26.1 を確認する

```bash
go version
# go1.26.1 でなければ https://go.dev/dl/ からインストール
```

---

## 2. 開発の進め方

### ワークフロー

```
スマホ: GitHub Issue を作成 or 確認
  ↓
PC: Claude Code を起動して「Issue #N を実装して」
  ↓
Claude Code がコード作成（サブスク内、追加料金なし）
  ↓
確認 → コミット・プッシュ
```

### 開始方法

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

# ビルド+起動（Windows では go run はファイアウォール警告が出るため使わない）
go build -o orangesql.exe . && ./orangesql.exe

# 別の DB ファイルを指定
go build -o orangesql.exe . && ./orangesql.exe -db /path/to/mydb.sqlite

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
cd OrangeSQL && go build -o orangesql.exe . && ./orangesql.exe

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
