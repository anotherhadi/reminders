package sidebar

import (
	"reminders/style"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	Label, Desc string
}

func (i Item) Title() string       { return i.Label }
func (i Item) Description() string { return i.Desc }
func (i Item) FilterValue() string { return i.Label }

func ChangeFocus() tea.Msg {
	return ChangeFocusMsg{}
}

type ChangeFocusMsg struct{}

type Model struct {
	List  list.Model
	Focus bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			if m.List.FilterState() != list.Filtering {
				return m, ChangeFocus
			}
		}
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width-2, msg.Height)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Focus {
		m.List.Styles.Title = style.TitleFocused
	} else {
		m.List.Styles.Title = style.TitleUnfocused
	}
	if len(m.List.Items()) == 0 {
		m.List.SetShowStatusBar(false)
	}
	return lipgloss.NewStyle().Width(m.List.Width()+2).Padding(0, 1).BorderStyle(lipgloss.RoundedBorder()).BorderForeground(style.MutedColor).BorderRight(true).Render(m.List.View())
}

func InitModel(title string, items []list.Item) Model {
	list := list.New(items, style.ListDelegate(), 0, 0)
	list.SetStatusBarItemName("list", "lists")
	list.Title = title
	list.SetShowHelp(false)
	list.Styles.Title = style.TitleFocused
	list.Styles.StatusBar = style.StatusBar
	return Model{
		List: list,
	}
}
