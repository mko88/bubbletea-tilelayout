package main

import (
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type DemoModel struct {
	layouts   []tl.TileLayout
	selected  int
	statusBar tl.Tile
}

func NewDemoModel() DemoModel {
	min := initialModelMinimal()
	weights := initialModelWeightsOnly()
	constraints := initialModelWithConstraints()
	many := initialModelManyLayouts()
	return DemoModel{
		layouts:  []tl.TileLayout{min, weights, constraints, many},
		selected: 0,
	}
}

func (d DemoModel) Init() tea.Cmd {
	return d.layouts[d.selected].Init()
}

func (d DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	_, cmd := d.layouts[d.selected].Update(msg)
	cmds = append(cmds, cmd)
	return d, tea.Batch(cmds...)
}

func (d DemoModel) updateSelection() (tea.Model, tea.Cmd) {
	if d.selected+1 < len(d.layouts) {
		d.selected += 1
	} else {
		d.selected = 0
	}
	return d, tea.WindowSize()
}

func (d DemoModel) View() string {
	return d.layouts[d.selected].View()
}
