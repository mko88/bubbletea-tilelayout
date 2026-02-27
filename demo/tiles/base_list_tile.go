package tiles

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	tl "github.com/mko88/bubbletea-tilelayout"
)

type BaseListTile struct {
	*tl.BaseTile
	Content   list.Model
	BoxBorder bool
}

func NewBaseListTile(items []list.Item, delegate list.ItemDelegate, size tl.Size, name string, boxBorder bool) BaseListTile {
	l := list.New(items, delegate, 10, 10)
	l.SetShowHelp(false)
	return BaseListTile{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: size,
		},
		Content:   l,
		BoxBorder: boxBorder,
	}
}

func (lt *BaseListTile) Init() tea.Cmd { return nil }

func (lt *BaseListTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.TileUpdatedMsg:
		if lt.GetName() != msg.Name {
			// only react to parent updates
			return lt, nil
		}
		newWidth := lt.Size.Width
		newHeight := lt.Size.Height
		if lt.BoxBorder {
			newWidth -= BOX_PAD
			newHeight -= BOX_PAD
		}
		lt.Content.SetWidth(newWidth)
		lt.Content.SetHeight(newHeight)
		return lt, nil
	}
	var cmd tea.Cmd
	lt.Content, cmd = lt.Content.Update(msg)
	return lt, cmd
}

func (lt *BaseListTile) View() string {
	if lt.BoxBorder {
		return RenderBox(lt.Name, lt.Content.View(), lt.Size)
	}
	return lt.Content.View()
}
