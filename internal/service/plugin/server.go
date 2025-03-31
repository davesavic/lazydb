package plugin

type DatabasePluginServer struct {
	Impl DatabasePlugin
}

func (d *DatabasePluginServer) Name(args any, resp *string) error {
	*resp = d.Impl.Name()
	return nil
}

func (d *DatabasePluginServer) Connect(conn string, resp *any) error {
	return d.Impl.Connect(conn)
}

func (d *DatabasePluginServer) Run(query string, resp *QueryResult) error {
	result, err := d.Impl.Run(query)
	if err != nil {
		return err
	}

	resp.Columns = result.Columns
	resp.Rows = result.Rows

	return nil
}

func (d *DatabasePluginServer) Close(args any, resp *any) error {
	return d.Impl.Close()
}
