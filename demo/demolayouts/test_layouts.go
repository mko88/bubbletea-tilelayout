package main

import (
	"github.com/charmbracelet/bubbles/viewport"
	tl "github.com/mko88/bubbletea-tilelayout"
	"github.com/mko88/bubbletea-tilelayout/demo/tiles"
)

func NewViewportTile(size tl.Size, name string, boxBorder bool) tiles.ViewportTile {
	vp := viewport.New(10, 10)
	return tiles.ViewportTile{
		Name:      name,
		Content:   vp,
		Size:      size,
		BoxBorder: boxBorder,
	}
}

func NewViewportTileMinimal(size tl.Size, name string, boxBorder bool) tiles.ViewportTileMinimal {
	vp := viewport.New(10, 10)
	return tiles.ViewportTileMinimal{
		Name:      name,
		Content:   vp,
		Size:      size,
		BoxBorder: boxBorder,
	}
}

func NewCustomStatusTile(size tl.Size, name string, content string) tiles.CustomTile {
	return tiles.CustomTile{
		Name:    name,
		Content: content,
		Size:    size,
	}
}

func initialModelMinimal() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	box := NewViewportTileMinimal(tl.Size{Weight: 1.00}, "Box1", true)
	status := NewCustomStatusTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	root.Add(&box)
	root.Add(&status)
	return root
}

func initialModelWithConstraints() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	sub1 := tl.TileLayout{
		Name:      "Sub-1",
		Direction: tl.Horizontal,
		Size:      tl.Size{Weight: 1.0},
	}
	status := NewCustomStatusTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	box1 := NewViewportTile(tl.Size{Weight: 0.20, FixedWidth: 50, FixedHeight: 15}, "Box1", true)
	box2 := NewViewportTile(tl.Size{Weight: 0.40, MaxWidth: 50}, "Box2", true)
	sub1.Add(&box1)
	sub1.Add(&box2)

	sub2 := tl.TileLayout{
		Name:      "Sub-2",
		Direction: tl.Vertical,
		Size:      tl.Size{Weight: 0.60},
	}
	box3 := NewViewportTile(tl.Size{Weight: 0.20, MinHeight: 6, MaxWidth: 90}, "Box3", true)
	box4 := NewViewportTile(tl.Size{Weight: 0.30, MaxWidth: 50}, "Box4", true)
	sub2.Add(&box3)
	sub2.Add(&box4)

	subsub1 := tl.TileLayout{
		Name:      "Sub-2-Sub-1",
		Direction: tl.Horizontal,
		Size:      tl.Size{Weight: 0.5},
	}
	box5 := NewViewportTile(tl.Size{Weight: 0.40, MaxHeight: 8}, "Box5", true)
	box6 := NewViewportTile(tl.Size{Weight: 0.60, MaxWidth: 40, MaxHeight: 14}, "Box6", true)
	subsub1.Add(&box5)
	subsub1.Add(&box6)
	sub2.Add(&subsub1)
	sub1.Add(&sub2)

	root.Add(&sub1)
	root.Add(&status)
	return root
}

func initialModelWeightsOnly() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	sub1 := tl.TileLayout{
		Name:      "Sub-1",
		Direction: tl.Horizontal,
		Size:      tl.Size{Weight: 1.0},
	}
	status := NewCustomStatusTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	box1 := NewViewportTile(tl.Size{Weight: 0.20}, "Box1", true)
	box2 := NewViewportTile(tl.Size{Weight: 0.40}, "Box2", true)
	sub1.Add(&box1)
	sub1.Add(&box2)

	sub2 := tl.TileLayout{
		Name:      "Sub-2",
		Direction: tl.Vertical,
		Size:      tl.Size{Weight: 0.60},
	}
	box3 := NewViewportTile(tl.Size{Weight: 0.20}, "Box3", true)
	box4 := NewViewportTile(tl.Size{Weight: 0.30}, "Box4", true)
	sub2.Add(&box3)
	sub2.Add(&box4)

	subsub1 := tl.TileLayout{
		Name:      "Sub-2-Sub-1",
		Direction: tl.Horizontal,
		Size:      tl.Size{Weight: 0.5},
	}
	box5 := NewViewportTile(tl.Size{Weight: 0.40}, "Box5", true)
	box6 := NewViewportTile(tl.Size{Weight: 0.60}, "Box6", true)
	subsub1.Add(&box5)
	subsub1.Add(&box6)
	sub2.Add(&subsub1)
	sub1.Add(&sub2)

	root.Add(&sub1)
	root.Add(&status)
	return root
}
