// SOURCE: https://github.com/charmbracelet/bubbles/blob/d3bd075ed2b27a3b5d76bb79b5d1c928dcd780d0/stopwatch/stopwatch.go
// Modified to allow setting the duration and adding a max duration

package stopwatch

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	lastID int
	idMtx  sync.Mutex
)

func nextID() int {
	idMtx.Lock()
	defer idMtx.Unlock()
	lastID++
	return lastID
}

// TickMsg is a message that is sent on every timer tick.
type TickMsg struct {
	// ID is the identifier of the stopwatch that sends the message. This makes
	// it possible to determine which stopwatch a tick belongs to when there
	// are multiple stopwatches running.
	//
	// Note, however, that a stopwatch will reject ticks from other
	// stopwatches, so it's safe to flow all TickMsgs through all stopwatches
	// and have them still behave appropriately.
	ID int
}

// StartStopMsg is sent when the stopwatch should start or stop.
type StartStopMsg struct {
	ID      int
	running bool
}

// SetMsg is sent when the stopwatch should set a new duration.
type SetMsg struct {
	ID int
	d  time.Duration
}

// Model for the stopwatch component.
type Model struct {
	d       time.Duration
	id      int
	running bool

	// How long to wait before every tick. Defaults to 1 second.
	Interval time.Duration

	// The max value the duration can be and will stop running once hit. A value of 0 means no max.
	Max time.Duration
}

// NewWithInterval creates a new stopwatch with the given timeout and tick
// interval.
func NewWithInterval(interval time.Duration) Model {
	return Model{
		Interval: interval,
		id:       nextID(),
	}
}

// New creates a new stopwatch with 1s interval.
func New() Model {
	return NewWithInterval(time.Second)
}

// ID returns the unique ID of the model.
func (m Model) ID() int {
	return m.id
}

// Init starts the stopwatch.
func (m Model) Init() tea.Cmd {
	return m.Start()
}

// Start starts the stopwatch.
func (m Model) Start() tea.Cmd {
	return tea.Batch(func() tea.Msg {
		return StartStopMsg{ID: m.id, running: true}
	}, tick(m.id, m.Interval))
}

// Stop stops the stopwatch.
func (m Model) Stop() tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{ID: m.id, running: false}
	}
}

// Toggle stops the stopwatch if it is running and starts it if it is stopped.
func (m Model) Toggle() tea.Cmd {
	if m.Running() {
		return m.Stop()
	}
	return m.Start()
}

// Set sets the stopwatch to the specified duration.
func (m Model) Set(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		return SetMsg{ID: m.id, d: d}
	}
}

// Reset resets the stopwatch to 0.
func (m Model) Reset() tea.Cmd {
	return m.Set(0)
}

// Running returns true if the stopwatch is running or false if it is stopped.
func (m Model) Running() bool {
	return m.running
}

// Update handles the timer tick.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StartStopMsg:
		if msg.ID != m.id {
			return m, nil
		}
		m.running = msg.running
		// Reset duration if starting from max
		if m.Max != 0 && m.running && m.d == m.Max {
			m.d = 0
		}
	case SetMsg:
		if msg.ID != m.id {
			return m, nil
		}
		m.d = msg.d
		if m.d < 0 {
			m.d = 0
		}
		if m.Max != 0 && m.d > m.Max {
			m.d = m.Max
			return m, m.Stop()
		}
	case TickMsg:
		if !m.running || msg.ID != m.id {
			break
		}
		m.d += m.Interval
		if m.Max != 0 && m.d > m.Max {
			m.d = m.Max
			return m, m.Stop()
		}
		return m, tick(m.id, m.Interval)
	}

	return m, nil
}

// Elapsed returns the time elapsed.
func (m Model) Elapsed() time.Duration {
	return m.d
}

// View of the timer component.
func (m Model) View() string {
	return m.d.String()
}

func tick(id int, d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return TickMsg{ID: id}
	})
}
