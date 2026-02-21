package tilelayout

import (
	"fmt"
	"time"

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
	BeforeLayout(tl *TileLayout)
	AfterLayout(tl *TileLayout)
}

type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

type Metrics struct {
	RenderTime  time.Duration
	RenderCount int
	TotalTime   time.Duration
	AverageTime time.Duration
}

// Optional: Get metrics report
func (tl *TileLayout) GetMetricsReport() string {
	return fmt.Sprintf(
		"Render Count: %d "+
			"Last Duration: %v "+
			"Average Duration: %v "+
			"Total Time: %v",
		tl.Metrics.RenderCount,
		tl.Metrics.RenderTime,
		tl.Metrics.AverageTime,
		tl.Metrics.TotalTime,
	)
}

type TileLayout struct {
	Name             string
	Size             Size
	Tiles            []Tile
	Proportions      []float64
	Direction        Direction
	Root             bool
	Parent           Tile
	TotalFixedWidth  int
	TotalFixedHeight int
	Metrics          Metrics
}

// AfterLayout implements [Tile].
func (*TileLayout) AfterLayout(tl *TileLayout) {
	// panic("unimplemented")
}

// BeforeLayout implements [Tile].
func (*TileLayout) BeforeLayout(tl *TileLayout) {
	// panic("unimplemented")
}

func NewRoot(direction Direction) *TileLayout {
	return &TileLayout{
		Name:      "Root",
		Direction: direction,
		Root:      true,
	}
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
	tl.TotalFixedWidth += tile.GetSize().FixedWidth
	tl.TotalFixedHeight += tile.GetSize().FixedHeight
	tl.Tiles = append(tl.Tiles, tile)
}

func (tl *TileLayout) isRoot() bool {
	return tl.GetParent() == nil
}

func (tl *TileLayout) DoLayout(width, height int) {
	if tl.isRoot() {
		tl.Size.Width = width
		tl.Size.Height = height
		tl.Size.Weight = 1
	}
	tl.BeforeLayout(tl)
	start := time.Now()
	tl.layout()
	elapsed := time.Since(start)
	for _, tile := range tl.Tiles {
		tile.AfterLayout(tl)
	}
	tl.AfterLayout(tl)
	tl.Metrics.RenderTime = elapsed
	tl.Metrics.RenderCount++
	tl.Metrics.TotalTime += elapsed
	tl.Metrics.AverageTime = tl.Metrics.TotalTime / time.Duration(tl.Metrics.RenderCount)
}

func (tl *TileLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tl.DoLayout(msg.Width, msg.Height)
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
	somethingResized := true
	sanityCheck := 0
	for somethingResized {
		sanityCheck++
		totalWidth, totalHeight, somethingResized = tl.distributeLeftover(totalWidth, totalHeight)
		if sanityCheck > 100 {
			panic("layout is unable to size itself for more than 100 iterations")
		}
	}
}

func CanGrowHeight(t Tile) bool {
	size := t.GetSize()
	return size.FixedHeight == 0 && (size.MaxHeight == 0 || size.Height < size.MaxHeight)
}

func CanGrowWidth(t Tile) bool {
	size := t.GetSize()
	return size.FixedWidth == 0 && (size.MaxWidth == 0 || size.Width < size.MaxWidth)
}

func (tl *TileLayout) distributeLeftover(totalWidth, totalHeight int) (int, int, bool) {
	leftoverWidth := tl.Size.Width - totalWidth
	leftoverHeight := tl.Size.Height - totalHeight
	if leftoverHeight == 0 && leftoverWidth == 0 {
		return totalWidth, totalHeight, false
	}
	somethingResized := false
	sumGrowableWeight := 0.0
	for _, tile := range tl.Tiles {
		weight := tile.GetSize().Weight
		switch tl.Direction {
		case Horizontal:
			if CanGrowWidth(tile) {
				sumGrowableWeight += weight
			}
		case Vertical:
			if CanGrowHeight(tile) {
				sumGrowableWeight += weight
			}
		}
	}

	for _, tile := range tl.Tiles {
		switch tl.Direction {
		case Horizontal:
			if CanGrowWidth(tile) && leftoverWidth > 0 {
				size := tile.GetSize()
				normalizedWeight := size.Weight / sumGrowableWeight
				toAdd := max(1, leftoverWidth*int(normalizedWeight))
				if size.MaxWidth > 0 && size.Width+toAdd > size.MaxWidth {
					toAdd = (size.MaxWidth - size.Width)
				}
				size.Width += toAdd
				tile.SetSize(size)
				somethingResized = true
				totalWidth += toAdd
				leftoverWidth -= toAdd
			}
		case Vertical:
			if CanGrowHeight(tile) && leftoverHeight > 0 {
				size := tile.GetSize()
				normalizedWeight := size.Weight / sumGrowableWeight
				toAdd := max(1, leftoverHeight*int(normalizedWeight))
				if size.MaxHeight > 0 && size.Height+toAdd > size.MaxHeight {
					toAdd = (size.MaxHeight - size.Height)
				}
				size.Height += toAdd
				tile.SetSize(size)
				somethingResized = true
				totalHeight += toAdd
				leftoverHeight -= toAdd
			}
		}
	}
	return totalWidth, totalHeight, somethingResized
}

func decideWidth(s Size, layout *TileLayout) int {
	availableWidth := layout.Size.Width
	if s.FixedWidth > 0 {
		// if fixed, return the minimum of the fixed or the available
		return min(availableWidth, s.FixedWidth)
	}
	if layout.Direction == Horizontal {
		availableWidth -= layout.TotalFixedWidth
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
	if layout.Direction == Vertical {
		availableHeight -= layout.TotalFixedHeight
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
