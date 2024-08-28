package style

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	AccentColor           = lipgloss.Color("#A594FD")
	AccentColorForeground = lipgloss.Color("#0B0B0B")
	AccentAltColor        = lipgloss.Color("#433A69")

	MutedColor           = lipgloss.Color("#3E3E3E")
	MutedColorForeground = lipgloss.Color("#FCFCFC")
	MutedAltColor        = lipgloss.Color("#373737")

	TitleFocused   = lipgloss.NewStyle().Foreground(AccentColorForeground).Background(AccentColor).Padding(0, 2).Bold(true)
	TitleUnfocused = lipgloss.NewStyle().Foreground(MutedColorForeground).Background(MutedColor).Padding(0, 2).Bold(true)
	StatusBar      = lipgloss.NewStyle().Foreground(MutedColor).PaddingBottom(1).PaddingLeft(2)
)

func ListDelegate() list.ItemDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.BorderForeground(AccentColor).Foreground(AccentColor)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.BorderForeground(AccentColor).Foreground(AccentAltColor)
	d.Styles.NormalTitle = d.Styles.NormalTitle.Foreground(MutedColor)
	d.Styles.NormalDesc = d.Styles.NormalTitle.Foreground(MutedAltColor)
	return d
}
