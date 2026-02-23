package tiles

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTileMinimal struct {
	*tl.BaseTile
	Content   viewport.Model
	BoxBorder bool
}

func NewViewportTileMinimal(size tl.Size, name string, boxBorder bool) ViewportTileMinimal {
	vp := viewport.New(10, 10)
	return ViewportTileMinimal{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: size,
		},
		Content:   vp,
		BoxBorder: boxBorder,
	}
}

func (vt *ViewportTileMinimal) Init() tea.Cmd { return nil }

func (vt *ViewportTileMinimal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.LayoutUpdatedMsg:
		if vt.Parent.GetName() != msg.Name {
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
		box := NewLabeledBox()
		box.BoxStyle = box.BoxStyle.
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 0, 0, 0).
			Width(vt.Size.Width - BOX_PAD).
			Height(vt.Size.Height - BOX_PAD)
		return box.Render(vt.Name, vt.Content.View(), vt.Size.Width-BOX_PAD)
	}
	return vt.Content.View()
}
