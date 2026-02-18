package tilelayout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Size struct {
	Width  int
	Height int
	Weight float64
}

func NewSize() Size {
	return Size{
		Width:  0,
		Height: 0,
		Weight: 0.0,
	}
}

type Tile interface {
	tea.Model
	GetName() string
	GetSize() Size
	SetSize(size Size)
	GetParent() Tile
	SetParent(tile Tile)
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

const (
	BOX_PAD = 2
)

type TileLayout struct {
	Name        string
	LayoutSize  Size
	Tiles       []Tile
	Proportions []float64
	Direction   Direction
	Root        bool
	Parent      Tile
}

func (tl *TileLayout) GetName() string {
	return tl.Name
}

func (tl *TileLayout) GetSize() Size {
	return tl.LayoutSize
}

func (tl *TileLayout) SetSize(size Size) {
	tl.LayoutSize = size
}

func (tl *TileLayout) GetParent() Tile {
	return tl.Parent
}

func (tl *TileLayout) SetParent(parent Tile) {
	tl.Parent = parent
}

func (tl *TileLayout) Init() tea.Cmd {
	return nil
}

func (tl *TileLayout) Add(tile Tile) {
	tile.SetParent(tl)
	tl.Tiles = append(tl.Tiles, tile)
}

func (tl *TileLayout) isRoot() bool {
	return tl.GetParent() == nil
}

func (tl *TileLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if tl.isRoot() {
			tl.LayoutSize.Width = msg.Width
			tl.LayoutSize.Height = msg.Height
			tl.LayoutSize.Weight = 1
		}
		tl.layout()
	}

	for i, tile := range tl.Tiles {
		if tile != nil {
			updated, cmd := tile.Update(msg)
			tl.Tiles[i] = updated.(Tile)
			cmds = append(cmds, cmd)
		}
	}

	return tl, tea.Batch(cmds...)
}

func (tl *TileLayout) View() string {
	if len(tl.Tiles) == 0 {
		return ""
	}
	var views []string

	for _, tile := range tl.Tiles {
		if tile == nil {
			continue
		}
		views = append(views, tile.View())
	}

	if tl.Direction == Horizontal {
		return lipgloss.JoinHorizontal(lipgloss.Top, views...)
	}
	return lipgloss.JoinVertical(lipgloss.Left, views...)
}

func (tl *TileLayout) layout() {
	if len(tl.Tiles) == 0 {
		return
	}
	totalHeight := 0
	totalWidth := 0
	for _, tile := range tl.Tiles {
		if tile == nil {
			continue
		}
		weight := tile.GetSize().Weight
		if tl.Direction == Horizontal {
			totalHeight = tl.LayoutSize.Height
			tileWidth := int(float64(tl.LayoutSize.Width) * weight)
			totalWidth += tileWidth
			if tileWidth > 0 {
				tile.SetSize(Size{Width: tileWidth, Height: tl.LayoutSize.Height, Weight: weight})
			}
		} else {
			totalWidth = tl.LayoutSize.Width
			tileHeight := int(float64(tl.LayoutSize.Height) * weight)
			totalHeight += tileHeight
			if tileHeight > 0 {
				tile.SetSize(Size{Width: tl.LayoutSize.Width, Height: tileHeight, Weight: weight})
			}
		}
	}
	leftoverHeight := tl.LayoutSize.Height - totalHeight
	leftoverWidth := tl.LayoutSize.Width - totalWidth
	lastTile := tl.Tiles[len(tl.Tiles)-1]
	lastTileSize := lastTile.GetSize()
	if leftoverHeight != 0 {
		lastTile.SetSize(Size{Width: lastTileSize.Width, Height: lastTileSize.Height + leftoverHeight, Weight: lastTileSize.Weight})
	}
	if leftoverWidth != 0 {
		lastTile.SetSize(Size{Width: lastTileSize.Width + leftoverWidth, Height: lastTileSize.Height, Weight: lastTileSize.Weight})
	}
}
