package database

import (
	"context"
	"fmt"
	"time"

	"github.com/davesavic/lazydb/internal/service/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var _ DatabaseIntegration = (*Postgres)(nil)

type Postgres struct {
	conn *pgx.Conn
}

// Name implements DatabasePlugin.
func (p *Postgres) Name() string {
	panic("unimplemented")
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

type Column struct {
	Name string
	Type uint32
}

type Result struct {
	Columns []Column
	Rows    []map[string]any
}

func (p *Postgres) Run(query string) (*QueryResult, error) {
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := rows.FieldDescriptions()

	result := &QueryResult{
		Columns: make([]string, len(columns)),
		Rows:    make([]map[string]any, 0),
	}

	for i, col := range columns {
		result.Columns[i] = col.Name
	}

	for rows.Next() {
		rowValues, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]any)

		for i, col := range columns {
			if rowValues[i] == nil {
				row[col.Name] = nil
				continue
			}

			// Switch case for types that need special handling
			switch v := rowValues[i].(type) {
			case time.Time:
				row[col.Name] = v.Format(time.RFC3339)
			case [16]uint8:
				row[col.Name] = uuid.UUID(v).String()
			default:
				row[col.Name] = v
			}
		}

		result.Rows = append(result.Rows, row)
	}

	return result, nil
}

func (p *Postgres) Close() error {
	return p.conn.Close(context.Background())
}
