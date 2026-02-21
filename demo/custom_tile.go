package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type CustomTile struct {
	Name    string
	Content string
	Size    tl.Size
	Parent  tl.Tile
}

// AfterLayout implements [tilelayout.Tile].
func (ct *CustomTile) AfterLayout(tl *tl.TileLayout) {
	//panic("unimplemented")
	ct.Content = tl.GetMetricsReport()
}

// BeforeLayout implements [tilelayout.Tile].
func (ct *CustomTile) BeforeLayout(tl *tl.TileLayout) {
	//panic("unimplemented")
}

func (ct *CustomTile) GetName() string {
	return ct.Name
}

func (ct *CustomTile) GetSize() tl.Size {
	return ct.Size
}

func (ct *CustomTile) SetSize(size tl.Size) {
	ct.Size = size
}

func (ct *CustomTile) Init() tea.Cmd {
	return nil
}

func (ct *CustomTile) GetParent() tl.Tile {
	return ct.Parent
}

func (ct *CustomTile) SetParent(parent tl.Tile) {
	ct.Parent = parent
}

func (ct *CustomTile) Update(tea.Msg) (tea.Model, tea.Cmd) {
	// the layout will calculate the size, based on the weight, nothing to do here
	return ct, nil
}

func (ct *CustomTile) View() string {
	return lipgloss.NewStyle().Width(ct.Size.Width).MaxHeight(ct.Size.Height).Render(ct.Content)
}
