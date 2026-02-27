package main

import (
	tl "github.com/mko88/bubbletea-tilelayout"
	"github.com/mko88/bubbletea-tilelayout/demo/tiles"
)

func interactiveModel() tl.TileLayout {
	root := tl.NewRoot(tl.Horizontal)
	targetLayout := initialModelWithConstraints()
	layoutList := tiles.NewLayoutTreeListTile(targetLayout, tl.Size{Weight: 1, MaxWidth: 30})
	root.Add(&layoutList)
	root.Add(&targetLayout)
	return root
}

func initialModelWithConstraints() tl.TileLayout {
	// create the root layout
	root := tl.NewTileLayout("Target", tl.Vertical, tl.Size{Weight: 1})
	border := true
	// create the tiles and sub-layouts
	contentArea := tl.NewTileLayout("ContentArea", tl.Horizontal, tl.Size{Weight: 1.0})
	status := tiles.NewTextTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	overview := tiles.NewLayoutOverviewTile(tl.Size{Weight: 0.3}, "Layout overview", border, &root)
	rightArea := tl.NewTileLayout("RightArea", tl.Vertical, tl.Size{Weight: 0.6})
	box3 := tiles.NewViewportTile(tl.Size{Weight: 0.20, MinHeight: 6, MaxWidth: 90}, "Box3", border)
	box4 := tiles.NewViewportTile(tl.Size{Weight: 0.30, MinWidth: 40, MaxWidth: 50}, "Box4", border)
	rightAreaSub := tl.NewTileLayout("RightAreaSubLayout", tl.Horizontal, tl.Size{Weight: 0.5})
	box5 := tiles.NewViewportTile(tl.Size{Weight: 0.40, MaxHeight: 8}, "Box5", border)
	box6 := tiles.NewViewportTile(tl.Size{Weight: 0.60, MaxWidth: 40, MaxHeight: 14}, "Box6", border)

	// add the tiles and sub-layouts to the layouts
	root.Add(&contentArea)
	root.Add(&status)
	contentArea.Add(&overview)
	contentArea.Add(&rightArea)
	rightArea.Add(&box3)
	rightArea.Add(&box4)
	rightArea.Add(&rightAreaSub)
	rightAreaSub.Add(&box5)
	rightAreaSub.Add(&box6)

	return root
}
