package main

import (
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type DemoModel struct {
	layout    *tl.TileLayout
	statusBar tl.Tile
}

func (d DemoModel) Init() tea.Cmd {
	d.layout.Init()
	return nil
}

func (d DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	d.layout.Update(msg)
	return d, nil
}

func (d DemoModel) View() string {
	return d.layout.View()
}
