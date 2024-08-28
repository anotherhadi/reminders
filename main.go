package main

import (
	"fmt"
	"os"
	"reminders/content"
	"reminders/markdown"
	"reminders/sidebar"
	"strconv"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	sidebarModel  sidebar.Model
	contentModels []content.Model
	help          help.Model
	keys          content.KeyMap
	size          Size
	toggleSidebar bool
	md            markdown.Md

	focus Focus
}

func (m Model) Init() tea.Cmd {
	return nil
}

type ReloadMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.updateSizes(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "tab":
			if m.toggleSidebar {
				m = m.ChangeFocus()
			}
			return m, nil

		case "ctrl+b": // Toggle the sidebar
			m.toggleSidebar = !m.toggleSidebar
			if m.toggleSidebar {
				m = m.FocusSidebar()
			} else {
				m = m.FocusContent()
			}
			m = m.updateSizes(m.size.screenWidth+2, m.size.screenHeight+2) // To match the padding
			return m, nil

		case "r": // Reload the markdown files
			m = m.reloadLists()
			return m, nil

		case "H": // Select the previous list, loop
			listLength := len(m.sidebarModel.List.Items())
			if listLength != 0 {
				newindex := 0
				if m.sidebarModel.List.Index()-1 < 0 {
					newindex = listLength - 1
				} else {
					newindex = m.sidebarModel.List.Index() - 1
				}
				m.sidebarModel.List.Select(newindex)
				m = m.FocusContent()
			}

		case "L": // Select the next list, loop
			listLength := len(m.sidebarModel.List.Items())
			if listLength != 0 {
				newindex := 0
				if m.sidebarModel.List.Index()+1 >= listLength {
					newindex = 0
				} else {
					newindex = m.sidebarModel.List.Index() + 1
				}
				m.sidebarModel.List.Select(newindex)
				m = m.FocusContent()
			}

		case m.keys.Help.Help().Key:
			m.help.ShowAll = !m.help.ShowAll
			m = m.updateSizes(m.size.screenWidth+2, m.size.screenHeight+2+lipgloss.Height(m.help.View(m.keys))) // To match the padding
			_, _ = m.Update(tea.WindowSizeMsg{Width: m.size.screenWidth + 2, Height: m.size.screenHeight + 2 + lipgloss.Height(m.help.View(m.keys))})
			return m, nil

		default:
			if m.focus == FocusSidebar {
				m.sidebarModel, cmd = m.sidebarModel.Update(msg)
				if cmd != nil {
					m = m.ChangeFocus()
					return m, nil
				}
				return m, cmd
			} else {
				m.contentModels[m.sidebarModel.List.Index()], cmd = m.contentModels[m.sidebarModel.List.Index()].Update(msg)
				return m, cmd
			}
		}
	}

	return m, cmd
}

func (m Model) View() string {
	var s string

	sidebarView := m.sidebarModel.View()
	contentView := ""
	if m.sidebarModel.List.Index() >= len(m.contentModels) {
		contentView = lipgloss.Place(m.size.contentWidth, m.size.contentHeight, lipgloss.Center, lipgloss.Center, "No content")
	} else {
		contentView = m.contentModels[m.sidebarModel.List.Index()].View()
	}

	if m.toggleSidebar {
		s = lipgloss.JoinHorizontal(lipgloss.Top, sidebarView, contentView)
	} else {
		s = contentView
	}
	helpView := m.help.View(m.keys)
	return lipgloss.NewStyle().Padding(1).Render(
		lipgloss.JoinVertical(lipgloss.Top, s, helpView),
	)
}

func (m Model) reloadLists() Model {
	var err error
	m.md, err = m.md.LoadMarkdown()
	if err != nil {
		return m
	}
	for i, item := range m.md.Lists {
		contentItems := []list.Item{}
		for _, i := range item.TodoItems {
			if i.Importance != "" {
				contentItems = append(contentItems, content.Item{Label: i.Importance + " " + i.State + " " + i.Label, Desc: i.Content})
			} else {
				contentItems = append(contentItems, content.Item{Label: i.State + " " + i.Label, Desc: i.Content})
			}
		}
		if len(m.contentModels) <= i {
			m.contentModels = append(m.contentModels, content.InitModel(item.Title, contentItems))
		} else {
			m.contentModels[i].List.SetItems(contentItems)
		}
	}
	return m
}

func run(filepath string, sidebarWith int, sidebarStatusOnStart bool) {
	md, err := markdown.Markdown(filepath)
	if err != nil {
		panic(err)
	}

	listItems := []list.Item{}
	for _, item := range md.Lists {
		var desc string
		if len(item.TodoItems) == 0 {
			desc = "No items"
		} else if len(item.TodoItems) == 1 {
			desc = "1 item"
		} else {
			desc = strconv.Itoa(len(item.TodoItems)) + " items"
		}

		listItems = append(listItems, sidebar.Item{Label: item.Title, Desc: desc})
	}

	sidebarModel := sidebar.InitModel(md.Title, listItems)

	m := Model{
		sidebarModel:  sidebarModel,
		toggleSidebar: sidebarStatusOnStart,
		size:          Size{defaultSidebarWidth: sidebarWith},
		md:            md,
		keys:          content.Keys,
		help:          help.New(),
	}
	m = m.reloadLists()
	m = m.initFocus()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func main() {
	run("./md_files/reminders.md", 30, false)
}
