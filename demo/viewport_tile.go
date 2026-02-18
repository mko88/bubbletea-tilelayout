package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTile struct {
	Name    string
	Box     LabeledBox
	Content viewport.Model
	Size    tl.Size
	Parent  tl.Tile
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

func (vt *ViewportTile) Update(tea.Msg) (tea.Model, tea.Cmd) {
	vt.Content.Width = vt.Size.Width - tl.BOX_PAD
	vt.Content.Height = vt.Size.Height - tl.BOX_PAD
	return vt, nil
}

func (vt *ViewportTile) View() string {
	vt.Box = NewLabeledBox()
	vt.Box.BoxStyle = vt.Box.BoxStyle.
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 0, 0, 0).
		Width(vt.Size.Width - tl.BOX_PAD).
		Height(vt.Size.Height - tl.BOX_PAD)
	s := fmt.Sprintf("%s (w:%d h:%d)", vt.Name, vt.Size.Width, vt.Size.Height)
	return vt.Box.Render(s, vt.Content.View(), vt.Size.Width-tl.BOX_PAD)

}
