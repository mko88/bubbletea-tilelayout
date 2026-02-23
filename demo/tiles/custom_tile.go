package tiles

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type CustomTile struct {
	*tl.BaseTile
	Content string
}

func NewCustomTile(size tl.Size, name string, content string) CustomTile {
	return CustomTile{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: size,
		},
		Content: content,
	}
}

func (ct *CustomTile) Init() tea.Cmd { return nil }

func (ct *CustomTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.LayoutUpdatedMsg:
		{
			if ct.Parent.GetName() != msg.Name {
				// only react to parent updates
				return ct, nil
			}
			ct.Content = fmt.Sprintf("%s: %v", msg.Name, msg.Metrics)
		}
	case tl.TileUpdatedMsg:
		if ct.GetName() == msg.Name {
			ct.Content = fmt.Sprintf("%v", msg)
		}
	}

	return ct, nil
}

func (ct *CustomTile) View() string {
	return lipgloss.NewStyle().Width(ct.Size.Width).MaxHeight(ct.Size.Height).Render(ct.Content)
}
