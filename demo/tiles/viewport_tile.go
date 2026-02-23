package tiles

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTile struct {
	*tl.BaseTile
	Content   viewport.Model
	BoxBorder bool
}

func NewViewportTile(size tl.Size, name string, boxBorder bool) ViewportTile {
	vp := viewport.New(10, 10)
	return ViewportTile{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: size,
		},
		Content:   vp,
		BoxBorder: boxBorder,
	}
}

func (vt *ViewportTile) Init() tea.Cmd { return nil }

func (vt *ViewportTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
		borderDescription := "I don't have a box border."
		if vt.BoxBorder {
			borderDescription = "I have a box border."
		}
		sizeDescription := fmt.Sprintf("Currently my dimensions are: %s.", printSize(vt.Size))
		text := fmt.Sprintf("I am viewport tile %v. %v %v\n%v", vt.Name, borderDescription, sizeDescription, msg.Metrics)
		text = lipgloss.NewStyle().Width(newWidth).Render(text)
		vt.Content.SetContent(text)
	}
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	vt.Content, cmd = vt.Content.Update(msg)
	cmds = append(cmds, cmd)
	return vt, tea.Batch(cmds...)
}

func (vt *ViewportTile) View() string {
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

var (
	zeroStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Dim gray
	positiveStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))  // Bright green
)

func styleInt(value int) string {
	if value > 0 {
		return positiveStyle.Render(fmt.Sprintf("%d", value))
	}
	return zeroStyle.Render(fmt.Sprintf("%d", value))
}

func styleFloat(value float64) string {
	if value > 0 {
		return positiveStyle.Render(fmt.Sprintf("%.2f", value))
	}
	return zeroStyle.Render(fmt.Sprintf("%.2f", value))
}

func printSize(s tl.Size) string {

	return fmt.Sprintf("actual[w:%s,h:%s,W:%s] min[w:%s,h:%s] max[w:%s,h:%s] fixed[w:%s,h:%s]",
		styleInt(s.Width), styleInt(s.Height), styleFloat(s.Weight),
		styleInt(s.MinWidth), styleInt(s.MinHeight),
		styleInt(s.MaxWidth), styleInt(s.MaxHeight),
		styleInt(s.FixedWidth), styleInt(s.FixedHeight))
}
