package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

func NewViewportTile(size tl.Size, name string, boxBorder bool) ViewportTile {
	vp := viewport.New(10, 10)
	return ViewportTile{
		Name:      name,
		Content:   vp,
		Size:      size,
		BoxBorder: boxBorder,
	}
}

func initialModel() *tl.TileLayout {
	layout := tl.TileLayout{
		Name:      "Root",
		Direction: tl.Horizontal,
		Root:      true,
	}

	box1 := NewViewportTile(tl.Size{Weight: 0.20, FixedWidth: 50, FixedHeight: 15}, "Box1", true)
	box2 := NewViewportTile(tl.Size{Weight: 0.40, MaxWidth: 50}, "Box2", true)
	layout.Add(&box1)
	layout.Add(&box2)

	sub1 := tl.TileLayout{
		Name:      "Sub-1",
		Direction: tl.Vertical,
		Size:      tl.Size{Weight: 0.60},
	}
	box3 := NewViewportTile(tl.Size{Weight: 0.20, MinHeight: 6}, "Box3", true)
	box4 := NewViewportTile(tl.Size{Weight: 0.30, MaxWidth: 50}, "Box4", true)
	sub1.Add(&box3)
	sub1.Add(&box4)

	subsub1 := tl.TileLayout{
		Name:      "SubSub-1",
		Direction: tl.Horizontal,
		Size:      tl.Size{Weight: 0.5},
	}
	box5 := NewViewportTile(tl.Size{Weight: 0.40, MaxHeight: 8}, "Box5", true)
	box6 := NewViewportTile(tl.Size{Weight: 0.60, MaxWidth: 14, MaxHeight: 14}, "Box6", true)
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
