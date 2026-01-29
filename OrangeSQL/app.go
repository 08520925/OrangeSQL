package main

import (
	"OrangeSQL/internal/app"
	"OrangeSQL/internal/domain"
	"OrangeSQL/internal/infra/sqlite"
	"context"
)

type App struct {
	ctx context.Context
	db  domain.Database
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 例：開発用に固定パス（あとで接続プロファイルにする）
	repo, err := sqlite.New("test.db")
	if err != nil {
		panic(err)
	}
	a.db = repo
}

// フロントから呼ぶ関数（Wailsがbinding生成する）
func (a *App) ExecuteSQL(sql string) (domain.QueryResult, error) {
	uc := app.ExecuteSQLUsecase{DB: a.db}
	return uc.Execute(sql)
}
