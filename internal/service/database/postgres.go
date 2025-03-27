package database

import (
	"context"
	"fmt"

	"github.com/davesavic/lazydb/internal/service/config"
	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	conn *pgx.Conn
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect(connCfg config.ConnectionConfig) error {
	conn, err := pgx.Connect(
		context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s", connCfg.User, connCfg.Password, connCfg.Host, connCfg.Port, connCfg.Database),
	)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	p.conn = conn

	err = p.conn.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	return nil
}

func (p *Postgres) GetTables() ([]string, error) {
	rows, err := p.conn.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		return nil, fmt.Errorf("could not get tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return nil, fmt.Errorf("could not scan table: %w", err)
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func (p *Postgres) Close() {
	p.conn.Close(context.Background())
}
