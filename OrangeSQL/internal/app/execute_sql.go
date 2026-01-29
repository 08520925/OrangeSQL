package app

import "OrangeSQL/internal/domain"

type ExecuteSQLUsecase struct {
	DB domain.Database
}

func (u ExecuteSQLUsecase) Execute(sql string) (domain.QueryResult, error) {
	return u.DB.Query(sql)
}
