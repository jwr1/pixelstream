package internal

import (
	"net/http"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode struct {
	currentFrame  *Frame
	appSwitchLock CmdLock
}

func NewViewMode() ViewMode {
	return ViewMode{
		currentFrame:  &Frame{},
		appSwitchLock: NewCmdLock(),
	}
}

func (m ViewMode) Init() tea.Cmd {
	return m.fetchFrame
}

func (m ViewMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			return NewMenuMode(), nil

		case "left":
			return m, m.appSwitchLock.TryLock(func() tea.Msg {
				http.Post(Host+"/api/previousapp", "application/json", http.NoBody)
				return nil
			})

		case "right":
			return m, m.appSwitchLock.TryLock(func() tea.Msg {
				http.Post(Host+"/api/nextapp", "application/json", http.NoBody)
				return nil
			})
		}

	case fetchFrameMsg:
		m.currentFrame = msg
		return m, m.fetchFrame
	}

	return m, nil
}

func (m ViewMode) View() string {
	var s strings.Builder

	s.WriteRune('\n')

	s.WriteString(m.currentFrame.View())

	s.WriteRune('\n')

	s.WriteString(helpStyle("[q] quit  [←] prev slide  [→] next slide\n"))

	return s.String()
}

type fetchFrameMsg *Frame

func (m ViewMode) fetchFrame() tea.Msg {
	m.currentFrame.ReceiveFrame(Host + "/api/screen")

	return fetchFrameMsg(m.currentFrame)
}
