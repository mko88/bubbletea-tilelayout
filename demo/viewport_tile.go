package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTile struct {
	Name      string
	Box       LabeledBox
	Content   viewport.Model
	Size      tl.Size
	Parent    tl.Tile
	BoxBorder bool
}

func (vt *ViewportTile) GetName() string {
	return vt.Name
}

func (vt *ViewportTile) GetSize() tl.Size {
	return vt.Size
}

func (vt *ViewportTile) SetSize(size tl.Size) {
	vt.Size = size
}

func (vt *ViewportTile) Init() tea.Cmd {
	return nil
}

func (vt *ViewportTile) GetParent() tl.Tile {
	return vt.Parent
}

func (vt *ViewportTile) SetParent(parent tl.Tile) {
	vt.Parent = parent
}

func (vt *ViewportTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.LayoutUpdatedMsg:
		if vt.Parent.GetName() != msg.Name {
			// only react to parent updates
			return vt, nil
		}
		newWidth := vt.Size.Width
		newHeight := vt.Size.Height
		if vt.BoxBorder {
			newWidth -= BOX_PAD
			newHeight -= BOX_PAD
		}
		vt.Content.Width = newWidth
		vt.Content.Height = newHeight
		borderDescription := "I don't have a box border."
		if vt.BoxBorder {
			borderDescription = "I have a box border."
		}
		sizeDescription := fmt.Sprintf("Currently my dimensions are: %s.", printSize(vt.Size))
		text := fmt.Sprintf("I am viewport tile %v. %v %v\n%v", vt.Name, borderDescription, sizeDescription, msg.Metrics)
		text = lipgloss.NewStyle().Width(newWidth).Render(text)
		vt.Content.SetContent(text)
	}
	return vt, nil
}

func (vt *ViewportTile) View() string {
	if vt.BoxBorder {
		vt.Box = NewLabeledBox()
		vt.Box.BoxStyle = vt.Box.BoxStyle.
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 0, 0, 0).
			Width(vt.Size.Width - BOX_PAD).
			Height(vt.Size.Height - BOX_PAD)
		return vt.Box.Render(vt.Name, vt.Content.View(), vt.Size.Width-BOX_PAD)
	}
	return vt.Content.View()
}

func printSize(s tl.Size) string {
	return fmt.Sprintf("actual[w:%d,h:%d,W:%.2f] min[w:%d,h:%d] max[w:%d,h:%d] fixed[w:%d,h:%d]",
		s.Width, s.Height, s.Weight, s.MinWidth, s.MinHeight, s.MaxWidth, s.MaxHeight, s.FixedWidth, s.FixedHeight)
}
