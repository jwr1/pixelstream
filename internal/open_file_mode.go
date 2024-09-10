package internal

import (
	"errors"
	"io/fs"
	"pixelstream/charmbracelet/bubbles/filepicker"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func NewOpenFileMode(fs fs.FS, initDir string) OpenFileMode {
	fp := filepicker.New(fs)
	fp.CurrentDirectory = initDir
	fp.ShowPermissions = false
	fp.AllowedTypes = []string{".pxlstrm", ".mp4", ".mkv", ".webm"}

	return OpenFileMode{
		filepicker: fp,
	}
}

type OpenFileMode struct {
	filepicker   filepicker.Model
	selectedFile string
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m OpenFileMode) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m OpenFileMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			return NewMenuMode(), nil
		}

	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		println(path)
		return SwitchMode(NewPlayMode(FileLocation{
			System: m.filepicker.FS,
			Path:   path,
		}))
	}

	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m OpenFileMode) View() string {
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else {
		s.WriteString(m.filepicker.Styles.Directory.Render(m.filepicker.CurrentDirectory))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}
