package tiles

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type LayoutTreeListTile struct {
	BaseListTile
	Layout tl.TileLayout
}

func NewLayoutTreeListTile(layout tl.TileLayout, size tl.Size) LayoutTreeListTile {
	rootModel := RowModel{
		Tile:   layout,
		Indent: 0,
	}
	models := make([]RowModel, 0)
	models = append(models, rootModel)
	layoutToTreeModel(layout, &models, 2)
	items := make([]list.Item, len(models))
	for i, m := range models {
		items[i] = m
	}

	delegate := &RowItemDelegate{}
	tree := NewBaseListTile(items, delegate, size, "Tree", true)
	return LayoutTreeListTile{
		BaseListTile: tree,
		Layout:       layout,
	}
}

func layoutToTreeModel(layout tl.TileLayout, models *[]RowModel, indent int) {
	for _, tile := range layout.Tiles {
		rm := RowModel{
			Tile:   tile,
			Indent: indent,
		}
		*models = append(*models, rm)
		if tile.IsLayout() {
			sub := tile.(*tl.TileLayout)
			layoutToTreeModel(*sub, models, indent+2)
		}

	}
}

// RowItemDelegate handles rendering of item rows
type RowItemDelegate struct{}

func (d RowItemDelegate) Height() int                             { return 1 }
func (d RowItemDelegate) Spacing() int                            { return 0 }
func (d RowItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d RowItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(RowModel)
	str := i.Tile.GetName()

	if index == m.Index() {
		str = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			Render("▶ " + strings.Repeat(" ", i.Indent) + str)
	} else {
		str = strings.Repeat(" ", i.Indent+2) + str
	}
	fmt.Fprint(w, str)
}

type RowModel struct {
	Tile   tl.Tile
	Indent int
}

func (rm RowModel) FilterValue() string { return rm.Tile.GetName() }
