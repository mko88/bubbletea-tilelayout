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
	l.SetShowFilter(false)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(false)
	return BaseListTile{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: size,
		},
		Content:   l,
		BoxBorder: boxBorder,
	}
}

func (blt *BaseListTile) Init() tea.Cmd { return nil }

func (blt *BaseListTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tl.TileUpdatedMsg:
		if blt.GetName() != msg.Name {
			// only react to parent updates
			return blt, nil
		}
		newWidth := blt.Size.Width
		newHeight := blt.Size.Height
		if blt.BoxBorder {
			newWidth -= BOX_PAD
			newHeight -= BOX_PAD
		}
		blt.Content.SetWidth(newWidth)
		blt.Content.SetHeight(newHeight)
		return blt, nil
	case tea.MouseMsg:
		return blt, nil
	}
	var cmd tea.Cmd
	blt.Content, cmd = blt.Content.Update(msg)
	return blt, cmd
}

func (blt *BaseListTile) View() string {
	if blt.BoxBorder {
		return RenderBox(blt.Name, blt.Content.View(), blt.Size)
	}
	return blt.Content.View()
}
