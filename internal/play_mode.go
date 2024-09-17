package internal

import (
	"errors"
	"io/fs"
	"pixelstream/charmbracelet/bubbles/stopwatch"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type playModeState int

const (
	playModeLoading playModeState = iota
	playModeConverting
	playModeError
	playModeReady
)

type PlayMode struct {
	state         playModeState
	stateMessage  string
	spinner       spinner.Model
	file          FileLocation
	pixelstream   *PixelStream
	frame         *Frame
	stopwatch     stopwatch.Model
	keymap        PlayModeKeymap
	help          help.Model
	progress      progress.Model
	sendFrameLock *CmdLock
}

type PlayModeKeymap struct {
	start         key.Binding
	stop          key.Binding
	reset         key.Binding
	quit          key.Binding
	skipBackwards key.Binding
	skipForwards  key.Binding
}

const defaultFrameRate = 16

func NewPlayMode(file FileLocation) PlayMode {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := PlayMode{
		state:   playModeLoading,
		spinner: s,
		file:    file,
		frame:   &Frame{},
		keymap: PlayModeKeymap{
			start: key.NewBinding(
				key.WithKeys(" ", "k"),
				key.WithHelp("space/k", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys(" ", "k"),
				key.WithHelp("space/k", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
			skipBackwards: key.NewBinding(
				key.WithKeys("left", "j"),
				key.WithHelp("←/j", "backwards"),
			),
			skipForwards: key.NewBinding(
				key.WithKeys("right", "l"),
				key.WithHelp("→/l", "forwards"),
			),
		},
		help:          help.New(),
		progress:      progress.New(progress.WithoutPercentage(), progress.WithWidth(46), progress.WithScaledGradient("#FF7CCB", "#FDFF8C")),
		sendFrameLock: &CmdLock{},
	}

	m.keymap.start.SetEnabled(false)

	return m
}

func (m PlayMode) Init() tea.Cmd {
	return m.LoadFile()
}

func (m PlayMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			return NewMenuMode(), nil
		}

		switch {
		case key.Matches(msg, m.keymap.reset):
			return m, m.stopwatch.Reset()
		case key.Matches(msg, m.keymap.start, m.keymap.stop):
			return m, m.stopwatch.Toggle()
		case key.Matches(msg, m.keymap.skipBackwards):
			return m, m.stopwatch.Set(m.stopwatch.Elapsed() - time.Second*5)
		case key.Matches(msg, m.keymap.skipForwards):
			return m, m.stopwatch.Set(m.stopwatch.Elapsed() + time.Second*5)
		}

	case playModeStateMsg:
		m.state = msg.state
		switch msg.state {
		case playModeLoading:
			return m, nil
		case playModeConverting:
			return m, m.GenerateFile()
		case playModeError:
			return m, nil
		case playModeReady:
			m.pixelstream = msg.pixelstream
			m.stopwatch = stopwatch.NewWithInterval(time.Second / time.Duration(m.pixelstream.FrameRate))
			m.stopwatch.Max = m.pixelstream.GetTotalDuration()
			return m, m.stopwatch.Init()
		}

	case stopwatch.StartStopMsg:
		m.keymap.stop.SetEnabled(!m.stopwatch.Running())
		m.keymap.start.SetEnabled(m.stopwatch.Running())
	case stopwatch.TickMsg:
		m.frame = m.pixelstream.GetFrame(m.stopwatch.Elapsed())

		cmd = m.sendFrameLock.TryLock(func() tea.Msg {
			m.frame.SendFrame(Host + "/api/notify")
			return nil
		})
	}

	var spinnerCmd tea.Cmd
	if m.state == playModeLoading || m.state == playModeConverting {
		m.spinner, spinnerCmd = m.spinner.Update(msg)
	}

	var stopwatchCmd tea.Cmd
	m.stopwatch, stopwatchCmd = m.stopwatch.Update(msg)

	return m, tea.Batch(cmd, spinnerCmd, stopwatchCmd)
}

func (m PlayMode) View() string {
	var s strings.Builder

	switch m.state {
	case playModeLoading:
		s.WriteString(m.spinner.View())
		s.WriteString("Loading file: ")
		s.WriteString(m.file.Path)
		s.WriteRune('\n')
	case playModeConverting:
		s.WriteString(m.spinner.View())
		s.WriteString("Converting file to .pxlstrm format...\n")
	case playModeError:
		s.WriteString("Error processing file: ")
		s.WriteString(m.file.Path)
		s.WriteRune('\n')
	case playModeReady:
		s.WriteString(m.frame.View())
		s.WriteRune('\n')

		s.WriteString(FmtDuration(m.stopwatch.Elapsed()))
		s.WriteRune(' ')
		s.WriteString(m.progress.ViewAs(float64(m.stopwatch.Elapsed()) / float64(m.stopwatch.Max)))
		s.WriteRune(' ')
		s.WriteString(FmtDuration(m.stopwatch.Max))

		s.WriteRune('\n')
	}

	if m.stateMessage != "" {
		s.WriteString(m.stateMessage)
		s.WriteRune('\n')
	}

	if m.state == playModeReady {
		s.WriteString(m.helpView())
	} else {
		s.WriteString(m.helpViewQuitOnly())
	}

	return s.String()
}

func (m PlayMode) helpViewQuitOnly() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.quit,
	})
}

func (m PlayMode) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
		m.keymap.skipBackwards,
		m.keymap.skipForwards,
	})
}

type playModeStateMsg struct {
	state        playModeState
	stateMessage string
	pixelstream  *PixelStream
}

func (m PlayMode) LoadFile() tea.Cmd {
	return tea.Sequence(
		m.spinner.Tick,
		func() tea.Msg {
			return playModeStateMsg{
				state: playModeLoading,
			}
		},
		func() tea.Msg {
			var pixelstream *PixelStream
			var err error
			if strings.HasSuffix(m.file.Path, pixelstreamFileExt) {
				pixelstream, err = LoadFile(m.file)
				if err != nil {
					return playModeStateMsg{
						state:        playModeError,
						stateMessage: err.Error(),
					}
				}
			} else if _, err := fs.Stat(m.file.System, m.file.Path+pixelstreamFileExt); !errors.Is(err, fs.ErrNotExist) {
				pixelstream, err = LoadFile(FileLocation{
					System: m.file.System,
					Path:   m.file.Path + pixelstreamFileExt,
				})
				if err != nil {
					return playModeStateMsg{
						state:        playModeError,
						stateMessage: err.Error(),
					}
				}
			} else {
				return playModeStateMsg{
					state: playModeConverting,
				}
			}

			return playModeStateMsg{
				state:       playModeReady,
				pixelstream: pixelstream,
			}
		},
	)
}

func (m PlayMode) GenerateFile() tea.Cmd {
	return func() tea.Msg {
		var pixelstream *PixelStream
		var err error
		pixelstream, err = GeneratePixelStream(m.file, defaultFrameRate)
		if err != nil {
			return playModeStateMsg{
				state:        playModeError,
				stateMessage: err.Error(),
			}
		}

		err = pixelstream.SaveFile(FileLocation{
			System: m.file.System,
			Path:   m.file.Path + pixelstreamFileExt,
		})
		if err != nil {
			return playModeStateMsg{
				state:        playModeError,
				stateMessage: err.Error(),
			}
		}

		return playModeStateMsg{
			state:       playModeReady,
			pixelstream: pixelstream,
		}
	}
}
