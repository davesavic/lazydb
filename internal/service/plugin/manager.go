package plugin

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	hcplugin "github.com/hashicorp/go-plugin"
)

type Manager struct {
	plugins map[string]*loadedPlugin
	logger  hclog.Logger
}

type loadedPlugin struct {
	name   string
	client *hcplugin.Client
	plugin DatabasePlugin
}

func NewManager(logger hclog.Logger) *Manager {
	return &Manager{
		plugins: make(map[string]*loadedPlugin),
		logger:  logger,
	}
}

func (m *Manager) GetPlugin(name string) (DatabasePlugin, bool) {
	loaded, ok := m.plugins[name]
	if !ok {
		return nil, false
	}

	return loaded.plugin, true
}

func (m *Manager) LoadPlugins(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		pluginPath := filepath.Join(path, file.Name())

		plug := hcplugin.NewClient(&hcplugin.ClientConfig{
			HandshakeConfig: hcplugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "BASIC_PLUGIN",
				MagicCookieValue: "hello",
			},
			Plugins: hcplugin.PluginSet{
				"database": &DatabasePluginAdaptor{},
			},
			Cmd:    exec.Command(pluginPath),
			Logger: m.logger,
		})

		rpcClient, err := plug.Client()
		if err != nil {
			plug.Kill()
			return err
		}

		raw, err := rpcClient.Dispense("database")
		if err != nil {
			plug.Kill()
			return err
		}

		dbPlugin := raw.(DatabasePlugin)

		m.plugins[file.Name()] = &loadedPlugin{
			name:   file.Name(),
			client: plug,
			plugin: dbPlugin,
		}
	}

	return nil
}

func (m *Manager) Close() {
	for _, plugin := range m.plugins {
		plugin.client.Kill()
	}
}
