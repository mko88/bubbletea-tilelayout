package tiles

import (
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTileMinimal struct {
	*BaseViewportTile
}

func NewViewportTileMinimal(size tl.Size, name string, boxBorder bool) ViewportTileMinimal {
	base := NewBaseViewportTile(size, name, boxBorder)
	return ViewportTileMinimal{
		BaseViewportTile: &base,
	}
}

func (vt *ViewportTileMinimal) Init() tea.Cmd { return nil }

func (vt *ViewportTileMinimal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.TileUpdatedMsg:
		if vt.Name != msg.Name {
			// only react to parent updates
			return vt, nil
		}
		newWidth := vt.Size.Width
		newHeight := vt.Size.Height
		if vt.BoxBorder {
			newWidth -= BOX_PAD
			newHeight -= BOX_PAD
		}
		vt.Content.Width = newWidth
		vt.Content.Height = newHeight
		vt.Content.SetContent("Press 'tab' to cycle to other layouts.\nTo quit press 'q' or 'ctrl+c'")
	}
	return vt, nil
}

func (vt *ViewportTileMinimal) View() string {
	if vt.BoxBorder {
		return RenderBox(vt.Name, vt.Content.View(), vt.Size)
	}
	return vt.Content.View()
}
