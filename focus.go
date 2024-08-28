package main

type Focus int

const (
	FocusSidebar Focus = 0
	FocusContent Focus = 1
)

func (m Model) initFocus() Model {
	if len(m.contentModels) > 0 {
		m.focus = FocusContent
		m.contentModels[0].Focus = true
	} else {
		m.focus = FocusSidebar
		m.sidebarModel.Focus = true
	}
	return m
}

func (m Model) ChangeFocus() Model {
	if m.focus == FocusSidebar {
		m = m.FocusContent()
	} else {
		m = m.FocusSidebar()
	}
	return m
}

func (m Model) FocusContent() Model {
	selected := m.sidebarModel.List.Index()
	if selected >= len(m.contentModels) {
		m = m.FocusSidebar()
		return m
	}
	m.focus = FocusContent
	for i := range m.contentModels {
		m.contentModels[i].Focus = true
	}
	m.sidebarModel.Focus = false
	return m
}

func (m Model) FocusSidebar() Model {
	selected := m.sidebarModel.List.Index()
	if selected >= len(m.contentModels) {
		m.sidebarModel.Focus = true
		return m
	}
	m.focus = FocusSidebar
	for i := range m.contentModels {
		m.contentModels[i].Focus = false
	}
	m.sidebarModel.Focus = true
	return m
}
