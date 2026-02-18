package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

func NewViewportTile(weight float64, name string) ViewportTile {
	vp := viewport.New(10, 10)
	vp.SetContent(fmt.Sprintf("%v\n------\nviewport content", name))
	return ViewportTile{
		Name:    name,
		Content: vp,
		Size:    tl.Size{Weight: weight},
	}
}

func initialModel() *tl.TileLayout {
	layout := tl.TileLayout{
		Name:      "root",
		Direction: tl.Horizontal,
		Root:      true,
	}

	box1 := NewViewportTile(0.20, "box1")
	box2 := NewViewportTile(0.40, "box2")
	layout.Add(&box1)
	layout.Add(&box2)

	sub1 := tl.TileLayout{
		Name:       "sub-1",
		Direction:  tl.Vertical,
		LayoutSize: tl.Size{Weight: 0.60},
	}
	box3 := NewViewportTile(0.20, "box3")
	box4 := NewViewportTile(0.30, "box4")
	sub1.Add(&box3)
	sub1.Add(&box4)

	subsub1 := tl.TileLayout{
		Name:       "subsub-1",
		Direction:  tl.Horizontal,
		LayoutSize: tl.Size{Weight: 0.5},
	}
	box5 := NewViewportTile(0.40, "box5")
	box6 := NewViewportTile(0.60, "box6")
	subsub1.Add(&box5)
	subsub1.Add(&box6)
	sub1.Add(&subsub1)

	layout.Add(&sub1)
	return &layout
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
