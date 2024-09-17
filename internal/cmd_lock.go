package internal

import (
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// Ensure a tea.Cmd can only run once at a time
type CmdLock struct {
	mutex sync.Mutex
}

func (m *CmdLock) Lock(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		m.mutex.Lock()
		msg := cmd()
		m.mutex.Unlock()
		return msg
	}
}

func (m *CmdLock) TryLock(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		if m.mutex.TryLock() {
			msg := cmd()
			m.mutex.Unlock()
			return msg
		}

		return nil
	}
}
