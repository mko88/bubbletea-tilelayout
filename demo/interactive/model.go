package main

import (
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type InteractiveModel struct {
	layouts   []tl.TileLayout
	selected  int
	statusBar tl.Tile
}

func NewInteractiveModel() InteractiveModel {
	m := interactiveModel()
	return InteractiveModel{
		layouts:  []tl.TileLayout{m},
		selected: 0,
	}
}

func (d InteractiveModel) Init() tea.Cmd {
	return d.layouts[d.selected].Init()
}

func (d InteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return d, tea.Quit
		case "tab":
			return d.updateSelection()
		}

	}
	updated, cmd := d.layouts[d.selected].Update(msg)
	d.layouts[d.selected] = updated.(tl.TileLayout)
	cmds = append(cmds, cmd)
	return d, tea.Batch(cmds...)
}

func (d InteractiveModel) updateSelection() (tea.Model, tea.Cmd) {

	if d.selected+1 < len(d.layouts) {
		d.selected += 1
	} else {
		d.selected = 0
	}
	return d, tea.WindowSize()
}

func (d InteractiveModel) View() string {
	return d.layouts[d.selected].View()
}
