package common

import (
	"github.com/davesavic/lazydb/internal/keybinding"
	"github.com/davesavic/lazydb/internal/message"
)

type ScreenProps struct {
	MessageManager *message.Manager
	Keymap         *keybinding.Keymap
}
