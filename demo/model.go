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
	return d.layout.Init()
}

func (d DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return d, tea.Quit
		}

	}
	_, cmd := d.layout.Update(msg)
	cmds = append(cmds, cmd)
	return d, tea.Batch(cmds...)
}

func (d DemoModel) View() string {
	return d.layout.View()
}
