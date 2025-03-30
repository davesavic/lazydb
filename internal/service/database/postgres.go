package database

import (
	"context"
	"fmt"

	"github.com/davesavic/lazydb/internal/service/config"
	"github.com/google/uuid"
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

type Column struct {
	Name string
	Type uint32
}

type Result struct {
	Columns []Column
	Rows    []map[string]any
}

func (p *Postgres) ExecuteQuery(query string) (*Result, error) {
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	fieldDescs := rows.FieldDescriptions()
	columns := make([]Column, len(fieldDescs))
	for i, fd := range fieldDescs {
		columns[i] = Column{
			Name: fd.Name,
			Type: fd.DataTypeOID,
		}
	}

	result := Result{
		Columns: columns,
	}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("could not get values: %w", err)
		}

		row := make(map[string]any)
		for i, value := range values {
			colName := columns[i].Name

			// UUID type's OID is 2950
			switch columns[i].Type {
			case 2950:
				if value != nil {
					row[colName] = uuid.UUID(value.([16]uint8)).String()
				}
			default:
				row[colName] = value
			}
		}

		result.Rows = append(result.Rows, row)
	}

	return &result, nil
}

func (p *Postgres) Close() {
	p.conn.Close(context.Background())
}
