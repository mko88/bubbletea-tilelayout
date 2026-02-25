package tiles

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type LayoutOverviewTile struct {
	BaseViewportTile
	Layout *tl.TileLayout
}

func NewLayoutOverviewTile(size tl.Size, name string, boxBorder bool, layout *tl.TileLayout) LayoutOverviewTile {
	base := NewBaseViewportTile(size, name, boxBorder)
	return LayoutOverviewTile{
		BaseViewportTile: base,
		Layout:           layout,
	}
}

func (lot *LayoutOverviewTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.TileUpdatedMsg:
		if msg.Name == lot.Name {
			lot.BaseViewportTile.Update(msg)
			var sb strings.Builder
			fmt.Fprintf(&sb, "---Layout tree---\n")
			printLayoutTree(&sb, *lot.Layout, "")
			fmt.Fprintf(&sb, "\n---Layout sizes---\n")
			printLayoutSizes(&sb, *lot.Layout)
			lot.Content.SetContent(sb.String())
			return lot, nil
		}
	}
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	lot.Content, cmd = lot.Content.Update(msg)
	cmds = append(cmds, cmd)
	return lot, tea.Batch(cmds...)
}

func printLayoutSizes(sb *strings.Builder, l tl.TileLayout) {
	for _, tile := range l.Tiles {
		if tl, ok := tile.(tl.TileLayout); ok {
			printLayoutSizes(sb, tl)
		}
	}
	name := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(l.Name)
	direction := ""
	switch l.Direction {
	case tl.Horizontal:
		direction = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Horizontal")
	case tl.Vertical:
		direction = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Vertical")
	}
	fmt.Fprintf(sb, "%v(%v)\n%v\n", name, direction, printSize(l.Size))

}

func printLayoutTree(sb *strings.Builder, l tl.TileLayout, prefix string) {
	name := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(l.Name)
	direction := ""
	switch l.Direction {
	case tl.Horizontal:
		direction = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Horizontal")
	case tl.Vertical:
		direction = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Vertical")
	}
	fmt.Fprintf(sb, "%s%v(%v)\n", prefix, name, direction)
	for _, tile := range l.Tiles {
		if tl, ok := tile.(tl.TileLayout); ok {
			printLayoutTree(sb, tl, prefix+"  ")
		} else {
			fmt.Fprintf(sb, "%s%v\n", prefix+"  ", tile.GetName())
		}
	}
}
