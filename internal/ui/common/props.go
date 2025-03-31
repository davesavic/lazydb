package common

import (
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/service/config"
	"github.com/davesavic/lazydb/internal/service/message"
	"github.com/davesavic/lazydb/internal/service/plugin"
)

type ScreenProps struct {
	MessageManager  *message.Manager
	ConfigService   *config.Service
	DatabaseService plugin.DatabasePlugin
	Keymap          *keybinding.Keymap
}
