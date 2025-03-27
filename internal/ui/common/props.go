package common

import (
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/service/config"
	"github.com/davesavic/lazydb/internal/service/database"
	"github.com/davesavic/lazydb/internal/service/message"
)

type ScreenProps struct {
	MessageManager  *message.Manager
	ConfigService   *config.Service
	DatabaseService *database.Postgres
	Keymap          *keybinding.Keymap
}
