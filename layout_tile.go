package tilelayout

import tea "github.com/charmbracelet/bubbletea"

type Tile interface {
	tea.Model
	GetName() string
	GetSize() Size
	SetSize(size Size)
	GetParent() Tile
	SetParent(tile Tile)
}

type BaseTile struct {
	Name   string
	Size   Size
	Parent Tile
}

func (bt BaseTile) GetName() string        { return bt.Name }
func (bt BaseTile) GetSize() Size          { return bt.Size }
func (bt *BaseTile) SetSize(size Size)     { bt.Size = size }
func (bt BaseTile) GetParent() Tile        { return bt.Parent }
func (bt *BaseTile) SetParent(parent Tile) { bt.Parent = parent }
