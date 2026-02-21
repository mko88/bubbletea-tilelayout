# Bubbletea TileLayout

A flexible and powerful tile-based layout system for [Bubble Tea](https://github.com/charmbracelet/bubbletea) terminal applications. Create complex, responsive terminal UI layouts with nested tiles, proportional sizing, and dynamic constraints.

## Features

- 🎯 **Flexible Layouts**: Create horizontal and vertical layouts that can be nested
- 📐 **Proportional Sizing**: Use weights to distribute space proportionally among tiles
- 🔒 **Size Constraints**: Set fixed, minimum, and maximum widths and heights for separate
- 🔄 **Responsive**: Automatically recalculates layouts on terminal resize
- 🧩 **Composable**: Implements the standard Bubble Tea `Model` interface for seamless integration
- 🎨 **Compatible**: Works with Bubble Tea components like viewports, lists etc. Custom tiles can also be implemented

## Installation

```bash
go get github.com/mko88/bubbletea-tilelayout
```

## QuickStart (Demo)

The `demo/` directory contains examples:

- **Minimal Layout**: Simple single-tile example
- **Weights Only**: Proportional sizing demonstration
- **Constraints**: Complex nested layouts with min/max/fixed constraints
- **Custom Tiles**: Examples of custom tile implementations

Run the demo:

```bash
cd demo
go build
go run .
```

## Core Concepts

### Tile Interface

All tiles must implement the `Tile` interface:

```go
type Tile interface {
    tea.Model
    GetName() string
    GetSize() Size
    SetSize(size Size)
    GetParent() Tile
    SetParent(tile Tile)
}
```

### Size Configuration

The `Size` struct provides flexible sizing options:

```go
type Size struct {
    Width       int     // Calculated width
    Height      int     // Calculated height
    Weight      float64 // Proportional weight (0.0 - 1.0)
    MinWidth    int     // Minimum width constraint
    MinHeight   int     // Minimum height constraint
    MaxWidth    int     // Maximum width constraint
    MaxHeight   int     // Maximum height constraint
    FixedWidth  int     // Fixed width (overrides weight)
    FixedHeight int     // Fixed height (overrides weight)
}
```

### Layout Directions

- `tl.Horizontal`: Arranges tiles side-by-side
- `tl.Vertical`: Stacks tiles top-to-bottom

### Messages
- `tl.LayoutUpdatedMsg`: Message sent when a layout is updated (layouted)

## Examples

### Basic Horizontal Layout

```go
root := tl.NewRoot(tl.Horizontal)
left := NewTile(tl.Size{Weight: 0.30})   // 30% width
right := NewTile(tl.Size{Weight: 0.70})  // 70% width
root.Add(&left)
root.Add(&right)
```

### Nested Layouts

```go
root := tl.NewRoot(tl.Vertical)

// Top section with horizontal split
top := &tl.TileLayout{
    Name:      "Top",
    Direction: tl.Horizontal,
    Size:      tl.Size{Weight: 0.80},
}
leftPane := NewTile(tl.Size{Weight: 0.50})
rightPane := NewTile(tl.Size{Weight: 0.50})
top.Add(&leftPane)
top.Add(&rightPane)

// Fixed height status bar at bottom
status := NewTile(tl.Size{FixedHeight: 1})

root.Add(top)
root.Add(&status)
```

### Using Constraints

```go
// Sidebar with fixed width and minimum height
sidebar := NewTile(tl.Size{
    FixedWidth: 30,
    MinHeight:  10,
})

// Main content with maximum width
main := NewTile(tl.Size{
    Weight:   0.70,
    MaxWidth: 100,
})

// Footer with fixed height
footer := NewTile(tl.Size{
    FixedHeight: 3,
})
```

## Creating Custom Tiles

Implement the `Tile` interface to create custom tiles:

```go
type MyTile struct {
    Name    string
    Size    tl.Size
    Parent  tl.Tile
    // Your custom fields
}

func (t *MyTile) GetName() string { return t.Name }
func (t *MyTile) GetSize() tl.Size { return t.Size }
func (t *MyTile) SetSize(size tl.Size) { t.Size = size }
func (t *MyTile) GetParent() tl.Tile { return t.Parent }
func (t *MyTile) SetParent(parent tl.Tile) { t.Parent = parent }

func (t *MyTile) Init() tea.Cmd {
    return nil
}

func (t *MyTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Handle updates
    return t, nil
}

func (t *MyTile) View() string {
    // Render using t.Size.Width and t.Size.Height
    return lipgloss.NewStyle().
        Width(t.Size.Width).
        Height(t.Size.Height).
        Render("My content")
}
```

## How It Works

1. **Initialization**: Create a root layout with `NewRoot(direction)`
2. **Composition**: Add tiles using `Add(tile)`, which can be layouts or custom tiles
3. **Auto Layout**: On `tea.WindowSizeMsg`, the layout automatically:
   - Calculates available space
   - Distributes space based on weights
   - Applies size constraints (min/max/fixed)
   - Updates all child tiles recursively
4. **Rendering**: Use `View()` to render tiles joined horizontally or vertically

## API Reference

### TileLayout

```go
// Create a new root layout
func NewRoot(direction Direction) *TileLayout

// Add a tile to the layout
func (tl *TileLayout) Add(tile Tile)

// Standard Bubble Tea methods
func (tl *TileLayout) Init() tea.Cmd
func (tl *TileLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (tl *TileLayout) View() string
```

### Size

```go
// Create a new Size struct with defaults
func NewSize() Size

// Create a copy of the Size
func (s *Size) Copy() Size
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style and layout rendering
- [Bubbles](https://github.com/charmbracelet/bubbles) - Common UI components (used in demos)

## License

Apache License 2.0

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Acknowledgments

Built with the excellent [Charm](https://charm.sh/) ecosystem of terminal UI tools.
