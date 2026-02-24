package main

import (
	tl "github.com/mko88/bubbletea-tilelayout"
	"github.com/mko88/bubbletea-tilelayout/demo/tiles"
)

func initialModelMinimal() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	box := tiles.NewViewportTileMinimal(tl.Size{Weight: 1.00}, "Box1", true)
	status := tiles.NewCustomTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	root.Add(&box)
	root.Add(&status)
	return root
}

func initialModelWithConstraints() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	sub1 := tl.NewTileLayout("Sub-1", tl.Horizontal, tl.Size{Weight: 1.0})
	status := tiles.NewCustomTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")

	box1 := tiles.NewViewportTile(tl.Size{Weight: 0.20, FixedWidth: 50, FixedHeight: 15}, "Box1", true)
	box2 := tiles.NewViewportTile(tl.Size{Weight: 0.40, MaxWidth: 50}, "Box2", true)
	sub1.Add(&box1)
	sub1.Add(&box2)

	sub2 := tl.NewTileLayout("Sub-2", tl.Vertical, tl.Size{Weight: 0.6})
	box3 := tiles.NewViewportTile(tl.Size{Weight: 0.20, MinHeight: 6, MaxWidth: 90}, "Box3", true)
	box4 := tiles.NewViewportTile(tl.Size{Weight: 0.30, MinWidth: 40, MaxWidth: 50}, "Box4", true)
	sub2.Add(&box3)
	sub2.Add(&box4)

	subsub1 := tl.NewTileLayout("Sub-2-Sub-1", tl.Horizontal, tl.Size{Weight: 0.5})
	box5 := tiles.NewViewportTile(tl.Size{Weight: 0.40, MaxHeight: 8}, "Box5", true)
	box6 := tiles.NewViewportTile(tl.Size{Weight: 0.60, MaxWidth: 40, MaxHeight: 14}, "Box6", true)
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
	sub1 := tl.NewTileLayout("Sub-1", tl.Horizontal, tl.Size{Weight: 1.0})
	status := tiles.NewCustomTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	box1 := tiles.NewViewportTile(tl.Size{Weight: 0.20}, "Box1", true)
	box2 := tiles.NewViewportTile(tl.Size{Weight: 0.40}, "Box2", true)
	sub1.Add(&box1)
	sub1.Add(&box2)

	sub2 := tl.NewTileLayout("Sub-2", tl.Vertical, tl.Size{Weight: 0.6})
	box3 := tiles.NewViewportTile(tl.Size{Weight: 0.20}, "Box3", true)
	box4 := tiles.NewViewportTile(tl.Size{Weight: 0.30}, "Box4", true)
	sub2.Add(&box3)
	sub2.Add(&box4)

	subsub1 := tl.NewTileLayout("Sub-2-Sub-1", tl.Horizontal, tl.Size{Weight: 0.5})
	box5 := tiles.NewViewportTile(tl.Size{Weight: 0.40}, "Box5", true)
	box6 := tiles.NewViewportTile(tl.Size{Weight: 0.60}, "Box6", true)
	subsub1.Add(&box5)
	subsub1.Add(&box6)
	sub2.Add(&subsub1)
	sub1.Add(&sub2)

	root.Add(&sub1)
	root.Add(&status)
	return root
}
