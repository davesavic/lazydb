package plugin

import (
	"net/rpc"

	hcplugin "github.com/hashicorp/go-plugin"
)

type QueryResult struct {
	Columns []string
	Rows    []map[string]any
}

type DatabasePlugin interface {
	Name() string
	Connect(conn string) error
	Run(query string) (*QueryResult, error)
	Close() error
}

type DatabasePluginAdaptor struct {
	Impl DatabasePlugin
}

// Client implements plugin.Plugin.
func (d *DatabasePluginAdaptor) Client(b *hcplugin.MuxBroker, c *rpc.Client) (any, error) {
	return &DatabasePluginClient{client: c}, nil
}

// Server implements plugin.Plugin.
func (d *DatabasePluginAdaptor) Server(*hcplugin.MuxBroker) (any, error) {
	return &DatabasePluginServer{Impl: d.Impl}, nil
}
