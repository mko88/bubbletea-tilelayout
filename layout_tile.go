package tilelayout

import tea "github.com/charmbracelet/bubbletea"

type Tile interface {
	tea.Model
	GetName() string
	GetSize() Size
	SetSize(size Size)
	GetParent() Tile
	SetParent(tile Tile)
	IsLayout() bool
}

type BaseTile struct {
	Name   string
	Size   Size
	Parent Tile
}

type TileUpdatedMsg struct {
	Name string
	Size Size
}

func NewTileUpdatedMsg(t Tile) tea.Cmd {
	return func() tea.Msg {
		return TileUpdatedMsg{
			Name: t.GetName(),
			Size: t.GetSize(),
		}
	}
}

func (bt BaseTile) GetName() string        { return bt.Name }
func (bt BaseTile) GetSize() Size          { return bt.Size }
func (bt *BaseTile) SetSize(size Size)     { bt.Size = size }
func (bt BaseTile) GetParent() Tile        { return bt.Parent }
func (bt *BaseTile) SetParent(parent Tile) { bt.Parent = parent }
func (vt BaseTile) IsLayout() bool         { return false }
