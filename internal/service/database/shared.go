package database

import "github.com/davesavic/lazydb/internal/service/config"

type QueryResult struct {
	Columns []string
	Rows    []map[string]any
}

type DatabaseIntegration interface {
	Name() string
	Connect(config.ConnectionConfig) error
	GetTables() ([]string, error)
	Run(query string) (*QueryResult, error)
	Close() error
}
