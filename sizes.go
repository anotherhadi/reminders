package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Size struct {
	defaultSidebarWidth         int
	screenWidth, screenHeight   int
	sidebarWidth, sidebarHeight int
	contentWidth, contentHeight int
}

func (m Model) updateSizes(screenWidth, screenHeight int) Model {
	m.size.screenWidth = screenWidth - 2                                          // leave room for the padding
	m.size.screenHeight = screenHeight - 2 - lipgloss.Height(m.help.View(m.keys)) // leave room for the padding and help bar
	m.size.sidebarHeight = m.size.screenHeight
	m.size.contentHeight = m.size.screenHeight
	if m.toggleSidebar {
		m.size.sidebarWidth = m.size.defaultSidebarWidth
		m.sidebarModel, _ = m.sidebarModel.Update(tea.WindowSizeMsg{Width: m.size.sidebarWidth, Height: m.size.sidebarHeight})
		m.size.contentWidth = m.size.screenWidth - lipgloss.Width(m.sidebarModel.View())
	} else {
		m.size.sidebarWidth = 0
		m.size.contentWidth = m.size.screenWidth
	}

	m.sidebarModel, _ = m.sidebarModel.Update(tea.WindowSizeMsg{Width: m.size.sidebarWidth, Height: m.size.sidebarHeight})
	for i := range m.contentModels {
		m.contentModels[i], _ = m.contentModels[i].Update(tea.WindowSizeMsg{Width: m.size.contentWidth, Height: m.size.contentHeight})
	}

	return m
}
