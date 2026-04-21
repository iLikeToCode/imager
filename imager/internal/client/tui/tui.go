package tui

import (
	"fmt"
	"os"

	"imager/internal/client/backend"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	pages []Page
	page  int
}

func initialModel(client *backend.Client) *model {
	m := &model{}
	m.pages = []Page{
		initialImageSelectPage(m, client),
	}
	return m
}

func (m *model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}

	return m.pages[m.page].Update(msg)
}

func (m *model) View() tea.View {
	return m.pages[m.page].View()
}

func RunTui(client *backend.Client) {
	// Source - https://stackoverflow.com/a/70060999
	// Posted by arctan2
	// Retrieved 2026-04-21, License - CC BY-SA 4.0
	os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})

	p := tea.NewProgram(initialModel(client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
