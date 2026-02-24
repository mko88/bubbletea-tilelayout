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
	case tl.TileUpdatedMsg:
		if vt.GetName() != msg.Name {
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
		sizeDescription := fmt.Sprintf("Currently my dimensions are: %s", printSize(vt.Size))
		text := fmt.Sprintf("I am viewport tile %v. %v %v", vt.Name, borderDescription, sizeDescription)
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
	gray  = lipgloss.NewStyle().Foreground(lipgloss.Color("240")) // Dim gray
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))  // Bright green
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("88"))  // Bright green
)

func grayInt(value int) string {
	return gray.Render(fmt.Sprintf("%d", value))
}

func greenInt(value int) string {
	return green.Render(fmt.Sprintf("%d", value))
}

func redInt(value int) string {
	return red.Render(fmt.Sprintf("%d", value))
}

func styleFloat(value float64) string {
	if value > 0 {
		return green.Render(fmt.Sprintf("%.2f", value))
	}
	return gray.Render(fmt.Sprintf("%.2f", value))
}

func printSize(s tl.Size) string {
	return fmt.Sprintf("\nreal  [w:%s,h:%s,W:%s]\nmin   [w:%s,h:%s]\nmax   [w:%s,h:%s]\nfixed [w:%s,h:%s]",
		style(s.Width, s.MinWidth, s.MaxWidth, s.FixedWidth), style(s.Height, s.MinHeight, s.MaxHeight, s.FixedHeight), styleFloat(s.Weight),
		style(s.MinWidth, s.Width, 0, 0), style(s.MinHeight, s.Height, 0, 0),
		style(s.MaxWidth, s.Width, 0, 0), style(s.MaxHeight, s.Height, 0, 0),
		style(s.FixedWidth, s.Width, 0, 0), style(s.FixedHeight, s.Height, 0, 0))
}

func style(target, constraint1, constraint2, constraint3 int) string {
	switch target {
	case 0:
		return grayInt(target)
	case constraint1, constraint2, constraint3:
		return redInt(target)
	}
	return greenInt(target)
}
