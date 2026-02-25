package tiles

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type TextTile struct {
	*tl.BaseTile
	Data    map[string]tl.Metrics
	Content string
}

func NewTextTile(size tl.Size, name string, content string) TextTile {
	return TextTile{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: size,
		},
		Data: make(map[string]tl.Metrics),
	}
}

func (ct *TextTile) Init() tea.Cmd { return nil }

func (ct *TextTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.LayoutUpdatedMsg:
		ct.Data[msg.Name] = msg.Metrics
		keys := make([]string, 0, len(ct.Data))
		for k := range ct.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%v", "Layouting times: ")
		for _, k := range keys {
			fmt.Fprintf(&sb, "%v[%v] ", k, ct.Data[k])
		}
		ct.Content = sb.String()
	case tl.TileUpdatedMsg:
		if ct.GetName() == msg.Name {
			//
		}
	}

	return ct, nil
}

func (ct *TextTile) View() string {
	return lipgloss.NewStyle().Width(ct.Size.Width).MaxHeight(ct.Size.Height).Render(ct.Content)
}
