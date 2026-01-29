package domain

type Database interface {
	Query(sql string) (QueryResult, error)
	Close() error
}
