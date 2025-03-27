package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
)

type ConnectionConfig struct {
	Type     string `toml:"type"`
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Database string `toml:"database"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type ConnectionsConfig struct {
	Connections map[string]ConnectionConfig `toml:"connections"`
}

type Service struct {
	ConnectionsConfig *ConnectionsConfig
}

func NewService() *Service {
	return &Service{
		ConnectionsConfig: &ConnectionsConfig{
			Connections: make(map[string]ConnectionConfig),
		},
	}
}

func (s *Service) GetConnection(name string) (*ConnectionConfig, error) {
	slog.Debug("s.ConnectionsConfig", "Connections", s.ConnectionsConfig.Connections)
	conn, ok := s.ConnectionsConfig.Connections[name]
	if !ok {
		return nil, fmt.Errorf("connection not found: %s", name)
	}

	return &conn, nil
}

func (s *Service) LoadConnections(path string) (*ConnectionsConfig, error) {
	var connections ConnectionsConfig

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("could not create file: %w", err)
		}
	}

	meta, err := toml.DecodeFile(path, &connections)
	if err != nil {
		return nil, fmt.Errorf("could not decode file: %w", err)
	}

	if len(meta.Undecoded()) > 0 {
		return nil, fmt.Errorf("could not decode file: %w", err)
	}

	if connections.Connections == nil {
		return nil, fmt.Errorf("no connections in file: %w", err)
	}

	s.ConnectionsConfig = &connections

	slog.Debug("config.LoadConnections", "s.ConnectionsConfig", s.ConnectionsConfig)

	return &connections, nil
}
