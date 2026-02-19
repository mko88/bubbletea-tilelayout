package tilelayout

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Size struct {
	Width       int
	Height      int
	Weight      float64
	MinWidth    int
	MinHeight   int
	MaxWidth    int
	MaxHeight   int
	FixedWidth  int
	FixedHeight int
}

func NewSize() Size {
	return Size{}
}

func (s *Size) Copy() Size {
	return Size{
		Width:       s.Width,
		Height:      s.Height,
		Weight:      s.Weight,
		MinWidth:    s.MinWidth,
		MinHeight:   s.MaxHeight,
		MaxWidth:    s.MaxWidth,
		MaxHeight:   s.MaxHeight,
		FixedWidth:  s.FixedWidth,
		FixedHeight: s.FixedHeight,
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

type TileLayout struct {
	Name        string
	Size        Size
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
	return tl.Size
}

func (tl *TileLayout) SetSize(size Size) {
	tl.Size = size
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
			tl.Size.Width = msg.Width
			tl.Size.Height = msg.Height
			tl.Size.Weight = 1
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
		newSize := tile.GetSize()
		switch tl.Direction {
		case Horizontal:
			// for horizontal, the total height is always concidered to be the layout height
			totalHeight = tl.Size.Height
			// decide the actual height in case of min/max/fixed
			tileHeight := decideHeight(newSize, tl)
			// in case the calculated width is more than the left available
			tileWidth := min(decideWidth(newSize, tl), tl.Size.Width-totalWidth)
			// append to total
			totalWidth += tileWidth
			if tileWidth > 0 {
				// set new size to tile
				newSize.Width = tileWidth
				newSize.Height = tileHeight
				tile.SetSize(newSize)
			}
		case Vertical:
			// for vertical, the total width is always concidered to be the layout width
			totalWidth = tl.Size.Width
			// in case the calculated height is more than the left available
			tileHeight := min(decideHeight(newSize, tl), tl.Size.Height-totalHeight)
			// decide the actual width in case of min/max/fixed
			tileWidth := decideWidth(newSize, tl)
			// append to total
			totalHeight += tileHeight
			if tileHeight > 0 {
				// set new size to tile
				newSize.Width = tileWidth
				newSize.Height = tileHeight
				tile.SetSize(newSize)
			}
		}
	}
}

func decideWidth(s Size, layout *TileLayout) int {
	availableWidth := layout.Size.Width
	if s.FixedWidth > 0 {
		// if fixed, return the minimum of the fixed or the available
		return min(availableWidth, s.FixedWidth)
	}

	// calculate height based on weight and available
	w := int(float64(availableWidth) * s.Weight)

	// if max defined and calculated is more, return the max
	if s.MaxWidth > 0 && w > s.MaxWidth {
		return s.MaxWidth
	}

	// if min defined and calculated is less, return the min
	if s.MinWidth > 0 && w < s.MinWidth {
		return s.MinWidth
	}

	if layout.Direction == Vertical {
		// for vertical layout, if min is defined, return the maximum of avalilable and minimum
		if s.MinWidth > 0 {
			return max(availableWidth, s.MinWidth)
		}
		// for vertical layout, if max is defined, return the minimum of available and maximum
		if s.MaxWidth > 0 {
			return min(availableWidth, s.MaxWidth)
		}
		// if none, return the available
		return availableWidth
	}
	// no constraints - return calculated
	return w
}

func decideHeight(s Size, layout *TileLayout) int {
	availableHeight := layout.Size.Height
	if s.FixedHeight > 0 {
		// if fixed, return the minimum of the fixed or the available
		return min(availableHeight, s.FixedHeight)
	}

	// calculate height based on weight and available
	h := int(float64(availableHeight) * s.Weight)

	// if max defined and calculated is more, return the max
	if s.MaxHeight > 0 && h > s.MaxHeight {
		return s.MaxHeight
	}

	// if min defined and calculated is less, return the min
	if s.MinHeight > 0 && h < s.MinHeight {
		return s.MinHeight
	}
	if layout.Direction == Horizontal {
		// for horizontal layout, if min is defined, return the maximum of avalilable and minimum
		if s.MinHeight > 0 {
			return max(availableHeight, s.MinHeight)
		}
		// for horizontal layout, if max is defined, return the minimum of available and maximum
		if s.MaxHeight > 0 {
			return min(availableHeight, s.MaxHeight)
		}
		// if none, return the available
		return availableHeight
	}
	// no constraints - return calculated
	return h
}
