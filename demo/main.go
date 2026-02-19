package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

func NewViewportTile(weight float64, name string, boxBorder bool) ViewportTile {
	vp := viewport.New(10, 10)
	return ViewportTile{
		Name:      name,
		Content:   vp,
		Size:      tl.Size{Weight: weight},
		BoxBorder: boxBorder,
	}
}

func initialModel() *tl.TileLayout {
	layout := tl.TileLayout{
		Name:      "Root",
		Direction: tl.Horizontal,
		Root:      true,
	}

	box1 := NewViewportTile(0.20, "Box1", true)
	box2 := NewViewportTile(0.40, "Box2", true)
	layout.Add(&box1)
	layout.Add(&box2)

	sub1 := tl.TileLayout{
		Name:      "Sub-1",
		Direction: tl.Vertical,
		Size:      tl.Size{Weight: 0.60},
	}
	box3 := NewViewportTile(0.20, "Box3", true)
	box4 := NewViewportTile(0.30, "Box4", false)
	sub1.Add(&box3)
	sub1.Add(&box4)

	subsub1 := tl.TileLayout{
		Name:      "SubSub-1",
		Direction: tl.Horizontal,
		Size:      tl.Size{Weight: 0.5},
	}
	box5 := NewViewportTile(0.40, "Box5", false)
	box6 := NewViewportTile(0.60, "Box6", true)
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
