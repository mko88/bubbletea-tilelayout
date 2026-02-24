package tiles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

const (
	BOX_PAD = 2
)

type LabeledBox struct {
	BoxStyle   lipgloss.Style
	LabelStyle lipgloss.Style
}

func NewLabeledBox() LabeledBox {
	return LabeledBox{
		BoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1),

		// You could, of course, also set background and foreground colors here
		// as well.
		LabelStyle: lipgloss.NewStyle().
			PaddingTop(0).
			PaddingBottom(0).
			PaddingLeft(1).
			PaddingRight(1),
	}
}

func RenderBox(label string, content string, size tl.Size) string {
	box := NewLabeledBox()
	box.BoxStyle = box.BoxStyle.
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 0, 0, 0).
		Width(size.Width - BOX_PAD).
		Height(size.Height - BOX_PAD)
	return box.Render(label, content, size.Width-BOX_PAD)
}

func (b LabeledBox) Render(label, content string, width int) string {
	var (
		// Query the box style for some of its border properties so we can
		// essentially take the top border apart and put it around the label.
		border          lipgloss.Border        = b.BoxStyle.GetBorderStyle()
		topBorderStyler func(...string) string = lipgloss.NewStyle().Foreground(b.BoxStyle.GetBorderTopForeground()).Render
		topLeft         string                 = topBorderStyler(border.TopLeft)
		topRight        string                 = topBorderStyler(border.TopRight)

		renderedLabel string = b.LabelStyle.Render(label)
	)

	// Render top row with the label
	borderWidth := b.BoxStyle.GetHorizontalBorderSize()
	cellsShort := max(0, width+borderWidth-lipgloss.Width(topLeft+topRight+renderedLabel))
	gap := strings.Repeat(border.Top, cellsShort)
	top := topLeft + renderedLabel + topBorderStyler(gap) + topRight

	// Render the rest of the box
	bottom := b.BoxStyle.
		BorderTop(false).
		Width(width).
		Render(content)

	// Stack the pieces
	return top + "\n" + bottom
}
