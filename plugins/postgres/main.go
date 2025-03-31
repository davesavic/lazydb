package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/davesavic/lazydb/internal/service/plugin"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/jackc/pgx/v5"
)

var _ plugin.DatabasePlugin = (*Postgres)(nil)

type Postgres struct {
	conn   *pgx.Conn
	logger hclog.Logger
}

// Close implements plugin.DatabasePlugin.
func (p *Postgres) Close() error {
	p.logger.Info("closing connection")
	return p.conn.Close(context.Background())
}

// Connect implements plugin.DatabasePlugin.
func (p *Postgres) Connect(conn string) error {
	p.logger.Info("connecting to database", "conn", conn)

	c, err := pgx.Connect(context.Background(), conn)
	if err != nil {
		return err
	}

	p.logger.Info("connected to database")

	p.conn = c

	p.logger.Info("pinging database")

	err = p.conn.Ping(context.Background())
	if err != nil {
		return err
	}

	p.logger.Info("pinged database")

	return nil
}

// Name implements plugin.DatabasePlugin.
func (p *Postgres) Name() string {
	p.logger.Info("getting name")
	return "postgres"
}

// Run implements plugin.DatabasePlugin.
func (p *Postgres) Run(query string) (*plugin.QueryResult, error) {
	p.logger.Info("running query", "query", query)
	rows, err := p.conn.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p.logger.Info("query executed")

	columns := rows.FieldDescriptions()
	p.logger.Info("got columns", "columns", columns)

	result := &plugin.QueryResult{
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

	p.logger.Info("got values", "values", result.Rows)

	return result, nil
}

func NewPostgres(logger hclog.Logger) *Postgres {
	logger.Info("creating new postgres plugin")
	return &Postgres{logger: logger}
}

func main() {
	logFile, err := os.OpenFile("postgres.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		os.Exit(1)
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "postgres",
		Level:      hclog.Trace,
		Output:     logFile,
		JSONFormat: true,
	})

	impl := NewPostgres(logger)

	pluginMap := map[string]hcplugin.Plugin{
		"database": &plugin.DatabasePluginAdaptor{Impl: impl},
	}

	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: hcplugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: "hello",
		},
		Plugins: pluginMap,
	})
}
