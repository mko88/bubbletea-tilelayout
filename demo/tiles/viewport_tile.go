package tiles

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type ViewportTile struct {
	Name      string
	Box       LabeledBox
	Content   viewport.Model
	Size      tl.Size
	Parent    tl.Tile
	BoxBorder bool
}

func (vt *ViewportTile) GetName() string          { return vt.Name }
func (vt *ViewportTile) GetSize() tl.Size         { return vt.Size }
func (vt *ViewportTile) SetSize(size tl.Size)     { vt.Size = size }
func (vt *ViewportTile) GetParent() tl.Tile       { return vt.Parent }
func (vt *ViewportTile) SetParent(parent tl.Tile) { vt.Parent = parent }
func (vt *ViewportTile) Init() tea.Cmd            { return nil }

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
	return vt, nil
}

func (vt *ViewportTile) View() string {
	if vt.BoxBorder {
		vt.Box = NewLabeledBox()
		vt.Box.BoxStyle = vt.Box.BoxStyle.
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 0, 0, 0).
			Width(vt.Size.Width - BOX_PAD).
			Height(vt.Size.Height - BOX_PAD)
		return vt.Box.Render(vt.Name, vt.Content.View(), vt.Size.Width-BOX_PAD)
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
