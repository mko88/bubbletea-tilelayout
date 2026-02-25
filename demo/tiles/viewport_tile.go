package tiles

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTile struct {
	*BaseViewportTile
}

func NewViewportTile(size tl.Size, name string, boxBorder bool) ViewportTile {
	base := NewBaseViewportTile(size, name, boxBorder)
	return ViewportTile{
		BaseViewportTile: &base,
	}
}

func (vt *ViewportTile) Init() tea.Cmd { return nil }

func (vt *ViewportTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tl.TileUpdatedMsg:
		if vt.Name != msg.Name {
			// only react to parent updates
			return vt, nil
		}
		vt.BaseViewportTile.Update(msg)
		parent := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62")).Render(vt.Parent.GetName())
		text := fmt.Sprintf("Parent: %v\n%v", parent, printSize(vt.Size))
		text = lipgloss.NewStyle().Width(vt.BaseViewportTile.Content.Width).Render(text)
		vt.Content.SetContent(text)
	}
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	vt.Content, cmd = vt.Content.Update(msg)
	cmds = append(cmds, cmd)
	return vt, tea.Batch(cmds...)
}
