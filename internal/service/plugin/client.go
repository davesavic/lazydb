package plugin

import (
	"net/rpc"
)

var _ DatabasePlugin = (*DatabasePluginClient)(nil)

type DatabasePluginClient struct {
	client *rpc.Client
}

// Close implements DatabasePlugin.
func (d *DatabasePluginClient) Close() error {
	err := d.client.Call("Plugin.Close", new(any), new(any))
	return err
}

// Connect implements DatabasePlugin.
func (d *DatabasePluginClient) Connect(conn string) error {
	err := d.client.Call("Plugin.Connect", conn, new(any))
	return err
}

// Name implements DatabasePlugin.
func (d *DatabasePluginClient) Name() string {
	var resp string
	err := d.client.Call("Plugin.Name", new(any), &resp)
	if err != nil {
		return ""
	}

	return resp
}

// Run implements DatabasePlugin.
func (d *DatabasePluginClient) Run(query string) (*QueryResult, error) {
	var resp *QueryResult
	err := d.client.Call("Plugin.Run", query, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
