package tilelayout

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// The Size structure and constraints for the layout
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

// Creates new Size
func NewSize() Size {
	return Size{}
}

// Message returned on layout update
type LayoutUpdatedMsg struct {
	Name    string
	Metrics string
}

// The command to return the LayoutUpdatedMsg
func (tl *TileLayout) layoutUpdated() tea.Cmd {
	return func() tea.Msg {
		return LayoutUpdatedMsg{
			Name:    tl.Name,
			Metrics: tl.GetMetricsReport(),
		}
	}
}

// The layout direction
type Direction int

const (
	Horizontal Direction = iota
	Vertical
)

// Metric values for the layout
type Metrics struct {
	Time        time.Duration
	RenderCount int
	TotalTime   time.Duration
	AverageTime time.Duration
}

// Return the layout metrics as string
func (tl *TileLayout) GetMetricsReport() string {
	return fmt.Sprintf(
		"LayoutName: %v; LayoutedCount: %d; LastLayoutDuration: %v; AverageDuration: %v; TotalTime: %v",
		tl.Name, tl.Metrics.RenderCount, tl.Metrics.Time, tl.Metrics.AverageTime, tl.Metrics.TotalTime,
	)
}

type TileLayout struct {
	*BaseTile
	Tiles            []Tile
	Direction        Direction
	TotalFixedWidth  int
	TotalFixedHeight int
	Metrics          Metrics
}

func NewRoot(direction Direction) TileLayout {
	return TileLayout{
		BaseTile: &BaseTile{
			Name: "Root",
		},
		Direction: direction,
	}
}

func NewTileLayout(name string, direction Direction, size Size) TileLayout {
	return TileLayout{
		BaseTile: &BaseTile{
			Name: name,
			Size: size,
		},
		Direction: direction,
	}
}

// Add a tile. The parent of the tile is set to the layout.
func (tl *TileLayout) Add(tile Tile) {
	tile.SetParent(tl)
	tl.TotalFixedWidth += tile.GetSize().FixedWidth
	tl.TotalFixedHeight += tile.GetSize().FixedHeight
	tl.Tiles = append(tl.Tiles, tile)
}

// If the layout have no parent, it's considered root.
func (tl *TileLayout) isRoot() bool {
	return tl.GetParent() == nil
}

// Handle the WindowSizeMsg
// If the layout is root, set its dimensions to the new window size and weight to 1.0.
// Proceeds with layouting itself and record its metrics.
func (tl *TileLayout) handleWindowSizeMsg(msg tea.WindowSizeMsg) {
	if tl.isRoot() {
		tl.Size.Width = msg.Width
		tl.Size.Height = msg.Height
		tl.Size.Weight = 1
	}
	start := time.Now()
	tl.layout()
	elapsed := time.Since(start)
	tl.Metrics.Time = elapsed
	tl.Metrics.RenderCount++
	tl.Metrics.TotalTime += elapsed
	tl.Metrics.AverageTime = tl.Metrics.TotalTime / time.Duration(tl.Metrics.RenderCount)
}

func (tl TileLayout) Init() tea.Cmd { return nil }

// Handle update messages from BubbleTea.
// On WidnowSizeMsg, the layout is "layouted" and LayoutUpdatedMsg is additionally returned.
// The message is forwarded to each tile.
func (tl TileLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tl.handleWindowSizeMsg(msg)
		cmds = append(cmds, tl.layoutUpdated())
	}

	for i, tile := range tl.Tiles {
		if tile != nil {
			updated, cmd := tile.Update(msg)
			tl.Tiles[i] = updated.(Tile)
			cmds = append(cmds, cmd)
		}
	}

	return &tl, tea.Batch(cmds...)
}

// Render all tiles, joining them together.
func (tl TileLayout) View() string {
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

// Perform dimension calculation for all tiles in the layout.
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
	// distribute the leftover spaces caused by constraints and rounding errors
	somethingResized := true
	sanityCheck := 0
	for somethingResized {
		// we may need to run it multiple times - some tiles may hit their constraints and
		// there would be still leftover space
		sanityCheck++
		totalWidth, totalHeight, somethingResized = tl.distributeLeftover(totalWidth, totalHeight)
		// but more than 100 times is quite unusual. Panic out.
		if sanityCheck > 100 {
			panic("layout is unable to size itself for more than 100 iterations")
		}
	}
}

// A tile can grow width when no fixed width is set and either no max width is set, or
// the current width is less than the max width
func canGrowWidth(t Tile) bool {
	size := t.GetSize()
	return size.FixedWidth == 0 && (size.MaxWidth == 0 || size.Width < size.MaxWidth)
}

// A tile can grow height when no fixed height is set and either no max height is set, or
// the current height is less than the max height
func canGrowHeight(t Tile) bool {
	size := t.GetSize()
	return size.FixedHeight == 0 && (size.MaxHeight == 0 || size.Height < size.MaxHeight)
}

// Distribute the leftover space. The sum weight of all growable tiles is calculated
// and distributed between them. Maximum constraints are respected. In case there is
// still leftover space, but no growable tiles, the dimensions are left as they are.
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
			if canGrowWidth(tile) {
				sumGrowableWeight += weight
			}
		case Vertical:
			if canGrowHeight(tile) {
				sumGrowableWeight += weight
			}
		}
	}

	for _, tile := range tl.Tiles {
		switch tl.Direction {
		case Horizontal:
			if canGrowWidth(tile) && leftoverWidth > 0 {
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
			if canGrowHeight(tile) && leftoverHeight > 0 {
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

// Decide the width of a tile in available space:
// 1. If fixed is defined, the minimum between the fixed and available is returned.
// 2. Total fixed width is subtracted from the available width.
// 3. Width is calucalated based on the weight.
// 4. If calculated is more than the max, max is returned.
// 5. If calculated is less than the min, min is returned.
// 6. If the layout is vertical, the available width would be used as default.
// 6.1 If min is defined, the maximum between the min and available is returned.
// 6.2 If max is defined, the minimum between the max and available is returned.
// 7. If there are no constraint, the calculated (3.) width is returned.
func decideWidth(s Size, layout *TileLayout) int {
	availableWidth := layout.Size.Width
	if s.FixedWidth > 0 {
		return min(availableWidth, s.FixedWidth)
	}
	if layout.Direction == Horizontal {
		availableWidth -= layout.TotalFixedWidth
	}
	w := int(float64(availableWidth) * s.Weight)
	if s.MaxWidth > 0 && w > s.MaxWidth {
		return s.MaxWidth
	}
	if s.MinWidth > 0 && w < s.MinWidth {
		return s.MinWidth
	}

	if layout.Direction == Vertical {
		if s.MinWidth > 0 {
			return max(availableWidth, s.MinWidth)
		}
		if s.MaxWidth > 0 {
			return min(availableWidth, s.MaxWidth)
		}
		return availableWidth
	}
	return w
}

// Decide the height of a tile in available space:
// 1. If fixed is defined, the minimum between the fixed and available is returned.
// 2. Total fixed height is subtracted from the available height.
// 3. Height is calucalated based on the weight.
// 4. If calculated is more than the max, max is returned.
// 5. If calculated is less than the min, min is returned.
// 6. If the layout is horizontal, the available height would be used as default.
// 6.1 If min is defined, the maximum between the min and available is returned.
// 6.2 If max is defined, the minimum between the max and available is returned.
// 7. If there are no constraint, the calculated (3.) height is returned.
func decideHeight(s Size, layout *TileLayout) int {
	availableHeight := layout.Size.Height
	if s.FixedHeight > 0 {
		return min(availableHeight, s.FixedHeight)
	}
	if layout.Direction == Vertical {
		availableHeight -= layout.TotalFixedHeight
	}
	h := int(float64(availableHeight) * s.Weight)
	if s.MaxHeight > 0 && h > s.MaxHeight {
		return s.MaxHeight
	}
	if s.MinHeight > 0 && h < s.MinHeight {
		return s.MinHeight
	}
	if layout.Direction == Horizontal {
		if s.MinHeight > 0 {
			return max(availableHeight, s.MinHeight)
		}
		if s.MaxHeight > 0 {
			return min(availableHeight, s.MaxHeight)
		}
		return availableHeight
	}
	return h
}
