package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
	"github.com/mko88/bubbletea-tilelayout/demo/tiles"
)

// func (c *ControlPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "up", "k":
// 			if c.selectedOption > 0 {
// 				c.selectedOption--
// 			}
// 		case "down", "j":
// 			if c.selectedOption < len(c.options)-1 {
// 				c.selectedOption++
// 			}
// 		case "enter", " ":
// 			// Return a custom message for the action
// 			return c, func() tea.Msg {
// 				return ActionMsg{Action: c.options[c.selectedOption]}
// 			}
// 		}
// 	}
// 	return c, nil
// }

// ActionMsg represents a control panel action
type ActionMsg struct {
	Action string
}

// Main DemoModel
type DemoModel struct {
	rootLayout    tl.TileLayout
	layoutTree    *tiles.LayoutTreeListTile
	statsTile     *tiles.LayoutOverviewTile
	contentLayout *tl.TileLayout
	idCounter     int
	width         int
	height        int
}

func initialModel() *DemoModel {
	// Create root layout (horizontal split between control/stats and content area)
	root := tl.NewRoot(tl.Horizontal)

	// Left panel (control + stats) - 30% width
	leftPanel := tl.NewTileLayout("LeftPanel", tl.Vertical, tl.Size{Weight: .3})

	// Right content area - 70% width, starts with a single tile
	rightArea := tl.NewTileLayout("ContentArea", tl.Vertical, tl.Size{Weight: .7})
	// Initial content tile
	content1 := tiles.NewViewportTile(tl.Size{Weight: 1.0}, "Tile-1", true)
	rightArea.Add(&content1)

	// Control panel at top of left side - 70% height
	layoutTree := tiles.NewLayoutTreeListTile(&rightArea, tl.Size{Weight: 0.50})

	// Stats at bottom of left side - 30% height
	statsTile := tiles.NewLayoutOverviewTile(tl.Size{Weight: 0.50}, "Root stats", true, &root)

	leftPanel.Add(&layoutTree)
	leftPanel.Add(&statsTile)

	root.Add(&leftPanel)
	root.Add(&rightArea)

	return &DemoModel{
		rootLayout:    root,
		layoutTree:    &layoutTree,
		statsTile:     &statsTile,
		contentLayout: &rightArea,
		idCounter:     1,
	}
}

func (m *DemoModel) Init() tea.Cmd {
	return nil
}

func (m *DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update the root layout with new size
		m.rootLayout.SetSize(tl.Size{
			Width:  msg.Width,
			Height: msg.Height,
		})

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		} else if msg.String() == "h" {
			return m, m.addSplit(tl.Horizontal)
		} else if msg.String() == "v" {
			return m, m.addSplit(tl.Vertical)
		}

	case ActionMsg:
		return m, m.handleAction(msg.Action)

	case tl.LayoutUpdatedMsg:
		// Force a rerender
		return m, nil
	}

	// Update child components
	var cmds []tea.Cmd

	// Update the root layout
	updatedRoot, cmd := m.rootLayout.Update(msg)
	if updatedRoot != nil {
		m.rootLayout = updatedRoot.(tl.TileLayout)
	}
	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *DemoModel) handleAction(action string) tea.Cmd {
	switch action {
	case "Add Horizontal Split":
		return m.addSplit(tl.Horizontal)
	case "Add Vertical Split":
		return m.addSplit(tl.Vertical)
	case "Remove Last":
		return m.removeLastTile()
	case "Reset Layout":
		return m.resetLayout()
	case "Quit":
		return tea.Quit
	}
	return nil
}

func (m *DemoModel) addSplit(direction tl.Direction) tea.Cmd {
	// Create a new layout to replace the current content
	if len(m.contentLayout.Tiles) > 0 {
		m.idCounter++
		newTile := tiles.NewViewportTile(tl.Size{Weight: 0.5}, fmt.Sprintf("Tile-%d", m.idCounter), true)

		// If current area is a single tile, convert to horizontal layout
		if l, ok := m.contentLayout.Tiles[0].(tl.TileLayout); ok {
			// Add to existing horizontal layout
			l.Add(&newTile)
			// Redistribute weights equally
			for i := range l.Tiles {
				l.Tiles[i].SetSize(tl.Size{Weight: 1.0 / float64(len(l.Tiles))})
			}
			// Replace the current tile with the new layout
			m.contentLayout.Tiles = []tl.Tile{l}
		} else {
			currentTile := m.contentLayout.Tiles[0]
			currentTile.SetSize(tl.Size{Weight: 0.5})

			// Create new horizontal layout
			hLayout := tl.NewTileLayout(fmt.Sprintf("%v-%d", direction.Name(), m.idCounter), direction, tl.Size{Weight: 1.0})

			hLayout.Add(currentTile)
			hLayout.Add(&newTile)

			// Replace the current tile with the new layout
			m.contentLayout.Tiles = []tl.Tile{hLayout}
		}
		m.rootLayout.Tiles[1] = m.contentLayout
	}

	return tea.WindowSize()
}

func (m *DemoModel) removeLastTile() tea.Cmd {
	// Find the right content area
	rightArea := m.rootLayout.Tiles[1].(tl.TileLayout)

	if len(rightArea.Tiles) > 1 {
		// Remove the last child
		rightArea.Tiles = rightArea.Tiles[:len(rightArea.Tiles)-1]

		// Redistribute weights equally among remaining children
		for i := range rightArea.Tiles {
			rightArea.Tiles[i].SetSize(tl.Size{Weight: 1.0 / float64(len(rightArea.Tiles))})
		}
	} else if len(rightArea.Tiles) == 1 {
		// If only one tile left, replace with a new content tile
		m.idCounter++
		//newTile := tiles.NewViewportTile(m.contentCount, "Reset content area", tl.Size{Weight: 1.0})
		//rightArea.Tiles = []tl.Tile{newTile}
		// newTile.SetParent(rightArea)
	}

	return tea.WindowSize()
}

func (m *DemoModel) resetLayout() tea.Cmd {
	// Reset to initial layout
	newModel := initialModel()
	m.rootLayout = newModel.rootLayout
	m.layoutTree = newModel.layoutTree
	m.statsTile = newModel.statsTile
	m.idCounter = newModel.idCounter

	return tea.WindowSize()
}

func (m *DemoModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	return m.rootLayout.View()
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
	}
}
