package tui

import (
	"fmt"

	"imager/gen/pb"
	"imager/internal/client/backend"

	tea "charm.land/bubbletea/v2"
)

type ImageSelectPage struct {
	m      *model
	images []*pb.Image
	cursor int
	chosen bool
}

func initialImageSelectPage(m *model, client *backend.Client) *ImageSelectPage {
	images, err := client.Images.ListImages()
	if err != nil {
		fmt.Println(err)
		return &ImageSelectPage{}
	}
	return &ImageSelectPage{
		m:      m,
		images: images,
		chosen: false,
	}
}

func (p *ImageSelectPage) Init() tea.Cmd {
	return nil
}

func (p *ImageSelectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if p.chosen {
				return p, nil
			}
			if p.cursor > 0 {
				p.cursor--
			}

		case "down", "j":
			if p.chosen {
				return p, nil
			}
			if p.cursor < len(p.images)-1 {
				p.cursor++
			}

		case "enter", "space":
			p.chosen = true
			p.m.page += 1
		}
	}

	return p.m, nil
}

func (p *ImageSelectPage) View() tea.View {
	if p.chosen {
		return tea.NewView(p.images[p.cursor].Name)
	}

	s := "Select an image:\n\n"

	for i, img := range p.images {

		cursor := " "
		if p.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, img.Name)
	}

	s += "\nPress enter to continue.\n"
	s += "\nPress q to quit.\n"

	return tea.NewView(s)
}
