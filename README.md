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

- **Layouts demo** (`demo/demolayouts/`): Four preset layouts (minimal, weights, constraints, many-tiles). Use `Tab` to cycle, `q` to quit.
- **Interactive demo** (`demo/interactive/`): Live editor — add and delete tiles, navigate selection, and watch the tree view update in real time.

```bash
cd demo/demolayouts && go run .
# or
cd demo/interactive && go run .
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
    IsLayout() bool
    IsFocused() bool
    SetFocused(bool)
}
```

### Size Configuration

The `Size` struct provides flexible sizing options:

```go
type Size struct {
    Width       int     // Calculated width
    Height      int     // Calculated height
    Weight      float64 // Proportional weight — relative to siblings, auto-normalized
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
- `tl.LayoutUpdatedMsg`: Sent when a layout recalculates its dimensions. Contains the layout name and render metrics.
- `tl.TileUpdatedMsg`: Sent after each resize for every tile in the tree. Contains the tile name and its new `Size`.

## Examples

### Basic Horizontal Layout

```go
root := tl.NewRoot(tl.Horizontal)
left := NewMyTile(tl.Size{Weight: 0.30})   // 30% width
right := NewMyTile(tl.Size{Weight: 0.70})  // 70% width
root.Add(left)
root.Add(right)
```

### Nested Layouts

```go
root := tl.NewRoot(tl.Vertical)

// Top section with horizontal split
top := tl.NewTileLayout("Top", tl.Horizontal, tl.Size{Weight: 0.80})
leftPane := NewMyTile(tl.Size{Weight: 0.50})
rightPane := NewMyTile(tl.Size{Weight: 0.50})
top.Add(leftPane)
top.Add(rightPane)

// Fixed height status bar at bottom
status := NewMyTile(tl.Size{FixedHeight: 1})

root.Add(top)
root.Add(status)
```

### Using Constraints

```go
// Sidebar with fixed width
sidebar := NewMyTile(tl.Size{
    FixedWidth: 30,
})

// Main content with maximum width
main := NewMyTile(tl.Size{
    Weight:   0.70,
    MaxWidth: 100,
})

// Footer with fixed height
footer := NewMyTile(tl.Size{
    FixedHeight: 3,
})
```

## Creating Custom Tiles

Embed `*tl.BaseTile` to satisfy the `Tile` interface boilerplate, then implement `Init()`, `Update()`, and `View()`:

```go
type MyTile struct {
    *tl.BaseTile
    // Your custom fields
}

func NewMyTile(size tl.Size, name string) *MyTile {
    return &MyTile{
        BaseTile: &tl.BaseTile{
            Name: name,
            Size: size,
        },
    }
}

func (t *MyTile) Init() tea.Cmd {
    return nil
}

func (t *MyTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tl.TileUpdatedMsg:
        if t.GetName() == msg.Name {
            // React to your tile being resized
        }
    }
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

1. **Initialization**: Create a root layout with `NewRoot(direction)`. Calling `Init()` on the root propagates to every tile in the tree.
2. **Composition**: Add tiles using `Add(tile)`, which can be layouts or custom tiles
3. **Auto Layout**: On `tea.WindowSizeMsg`, the layout automatically:
   - Calculates available space
   - Distributes space based on weights (auto-normalized — `1/1/1` and `0.33/0.33/0.33` produce the same result)
   - Applies size constraints (min/max/fixed)
   - Updates all child tiles recursively via `SetSize` and `TileUpdatedMsg`
4. **Rendering**: Use `View()` to render tiles joined horizontally or vertically

## API Reference

### TileLayout

```go
// Create a new root layout (no parent, fills the terminal window)
func NewRoot(direction Direction) *TileLayout

// Create a named sub-layout with a size configuration
func NewTileLayout(name string, direction Direction, size Size) *TileLayout

// Add a tile; sets the tile's parent
func (tl *TileLayout) Add(tile Tile)

// Replace a tile by name; returns true if found
func (tl *TileLayout) Replace(name string, newTile Tile) bool

// Remove a tile by name; returns true if found
func (tl *TileLayout) Remove(name string) bool

// Focus navigation — operates on all leaf tiles in tree order
func (tl *TileLayout) FocusFirst()
func (tl *TileLayout) FocusNext()
func (tl *TileLayout) FocusPrev()
func (tl *TileLayout) FocusedTile() Tile  // nil if nothing focused

// Standard Bubble Tea methods; Init() propagates to all child tiles
func (tl *TileLayout) Init() tea.Cmd
func (tl *TileLayout) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (tl *TileLayout) View() string
```

### BaseTile

```go
// Embed in your custom tile to satisfy the Tile interface
type BaseTile struct {
    Name    string
    Size    Size
    Parent  Tile
    Focused bool
}
```

### Size

```go
// Create a new Size struct with defaults
func NewSize() Size
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
