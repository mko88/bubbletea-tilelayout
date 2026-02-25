package main

import (
	tl "github.com/mko88/bubbletea-tilelayout"
	"github.com/mko88/bubbletea-tilelayout/demo/tiles"
)

func initialModelMinimal() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	box := tiles.NewViewportTileMinimal(tl.Size{Weight: 1.00}, "Box1", true)
	status := tiles.NewTextTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	root.Add(&box)
	root.Add(&status)
	return root
}

func initialModelWithConstraints() tl.TileLayout {
	// create the root layout
	root := tl.NewRoot(tl.Vertical)
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

func initialModelWeightsOnly() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	sub1 := tl.NewTileLayout("Sub-1", tl.Horizontal, tl.Size{Weight: 1.0})
	status := tiles.NewTextTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")
	box1 := tiles.NewViewportTile(tl.Size{Weight: 0.33}, "Box1", true)
	box2 := tiles.NewViewportTile(tl.Size{Weight: 0.33}, "Box2", true)
	sub1.Add(&box1)
	sub1.Add(&box2)

	sub2 := tl.NewTileLayout("Sub-2", tl.Vertical, tl.Size{Weight: 0.33})
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

func initialModelManyLayouts() tl.TileLayout {
	root := tl.NewRoot(tl.Vertical)
	content := tl.NewTileLayout("Content", tl.Horizontal, tl.Size{Weight: 1.0})

	left := tl.NewTileLayout("Left", tl.Vertical, tl.Size{Weight: .33})
	leftTop := tl.NewTileLayout("LeftTop", tl.Horizontal, tl.Size{Weight: .33})
	leftBottom := tl.NewTileLayout("LeftBottom", tl.Horizontal, tl.Size{Weight: .33})

	middle := tl.NewTileLayout("Middle", tl.Vertical, tl.Size{Weight: .33})

	right := tl.NewTileLayout("Right", tl.Horizontal, tl.Size{Weight: .33})

	ltb1 := tiles.NewLayoutOverviewTile(tl.Size{Weight: 0.6}, "lt-b1", true, &root)
	ltb2 := tiles.NewViewportTile(tl.Size{Weight: 0.4}, "lt-b2", true)

	lbb1 := tiles.NewViewportTile(tl.Size{Weight: 0.3}, "lb-b1", true)
	lbb2 := tiles.NewViewportTile(tl.Size{Weight: 0.7}, "lb-b2", true)

	mb1 := tiles.NewViewportTile(tl.Size{Weight: 0.33}, "m-b1", true)
	mb2 := tiles.NewViewportTile(tl.Size{Weight: 0.33}, "m-b2", true)
	mb3 := tiles.NewViewportTile(tl.Size{Weight: 0.33}, "m-b3", true)

	rb1 := tiles.NewViewportTile(tl.Size{Weight: 0.25}, "r-b1", true)
	rb2 := tiles.NewViewportTile(tl.Size{Weight: 0.25}, "r-b2", true)
	rb3 := tiles.NewViewportTile(tl.Size{Weight: 0.25}, "r-b3", true)
	rb4 := tiles.NewViewportTile(tl.Size{Weight: 0.25}, "r-b4", true)

	status := tiles.NewTextTile(tl.Size{FixedHeight: 1}, "Status", "I am the status tile. I have a fixed height of 1 and take up 100% space.")

	leftTop.Add(&ltb1)
	leftTop.Add(&ltb2)
	leftBottom.Add(&lbb1)
	leftBottom.Add(&lbb2)

	left.Add(&leftTop)
	left.Add(&leftBottom)

	middle.Add(&mb1)
	middle.Add(&mb2)
	middle.Add(&mb3)

	right.Add(&rb1)
	right.Add(&rb2)
	right.Add(&rb3)
	right.Add(&rb4)

	content.Add(&left)
	content.Add(&middle)
	content.Add(&right)

	root.Add(&content)
	root.Add(&status)
	return root
}
