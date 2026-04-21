package tui

import tea "charm.land/bubbletea/v2"

type Page interface {
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
	Init() tea.Cmd
}
