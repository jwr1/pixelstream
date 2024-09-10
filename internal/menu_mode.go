package internal

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	listHelpStyle     = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type MenuItem struct {
	Label string
	Mode  tea.Model
}

var MenuItems []MenuItem

func (i MenuItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(MenuItem)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(i.Label))
}

type MenuMode struct {
	list list.Model
}

func NewMenuMode() MenuMode {
	listItems := make([]list.Item, len(MenuItems))
	for i, v := range MenuItems {
		listItems[i] = v
	}

	l := list.New(listItems, itemDelegate{}, 80, 7)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = listHelpStyle

	l.KeyMap.ShowFullHelp.SetEnabled(false)

	return MenuMode{
		list: l,
	}
}

func (m MenuMode) Init() tea.Cmd {
	return nil
}

func (m MenuMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(MenuItem)
			if ok {
				return SwitchMode(i.Mode)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuMode) View() string {
	var s strings.Builder

	s.WriteString("pixelstream - Stream videos to your awtrix clock with ease.\n")
	s.WriteString("Host: ")
	s.WriteString(Host)
	s.WriteString("\n\n")

	s.WriteString(m.list.View())

	return s.String()
}
