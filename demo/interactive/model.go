package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tl "github.com/mko88/bubbletea-tilelayout"
)

// normalizeLayoutWeights sets equal weights for all tiles within a layout.
func normalizeLayoutWeights(l *tl.TileLayout) {
	n := float64(len(l.Tiles))
	if n == 0 {
		return
	}
	w := 1.0 / n
	for _, tile := range l.Tiles {
		size := tile.GetSize()
		size.Weight = w
		tile.SetSize(size)
	}
}

// collectNodes appends all nodes (layouts + leaves) under l to nodes, pre-order.
func collectNodes(l *tl.TileLayout, nodes *[]tl.Tile) {
	*nodes = append(*nodes, l)
	for _, tile := range l.Tiles {
		if child, ok := tile.(*tl.TileLayout); ok {
			collectNodes(child, nodes)
		} else {
			*nodes = append(*nodes, tile)
		}
	}
}

// ─── Preview tile ─────────────────────────────────────────────────────────────

type PreviewTile struct {
	*tl.BaseTile
}

func NewPreviewTile(name string) *PreviewTile {
	return &PreviewTile{
		BaseTile: &tl.BaseTile{
			Name: name,
			Size: tl.Size{Weight: 1.0},
		},
	}
}

func (pt *PreviewTile) Init() tea.Cmd { return nil }

func (pt *PreviewTile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tl.TileUpdatedMsg); ok && msg.Name == pt.Name {
		pt.Size = msg.Size
	}
	return pt, nil
}

func (pt *PreviewTile) View() string {
	w, h := pt.Size.Width, pt.Size.Height
	if w <= 2 || h <= 2 {
		return lipgloss.NewStyle().Width(w).Height(h).Render("")
	}
	borderColor := lipgloss.Color("63")
	if pt.IsFocused() {
		borderColor = lipgloss.Color("214")
	}
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	content := nameStyle.Render(pt.Name) + "\n" +
		dimStyle.Render(fmt.Sprintf("w:%-3d h:%-3d  wgt:%.2f", w, h, pt.Size.Weight))
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(w - 2).
		Height(h - 2).
		Render(content)
}

// ─── Control panel ────────────────────────────────────────────────────────────

const controlWidth = 38

type ControlPanel struct {
	*tl.BaseTile
	model *InteractiveModel
}

func NewControlPanel(model *InteractiveModel) *ControlPanel {
	return &ControlPanel{
		BaseTile: &tl.BaseTile{
			Name: "Controls",
			Size: tl.Size{FixedWidth: controlWidth},
		},
		model: model,
	}
}

func (cp *ControlPanel) Init() tea.Cmd { return nil }

func (cp *ControlPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tl.TileUpdatedMsg); ok && msg.Name == cp.Name {
		cp.Size = msg.Size
	}
	return cp, nil
}

func (cp *ControlPanel) View() string {
	m := cp.model
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	headStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	fmt.Fprintf(&sb, "%s\n\n", titleStyle.Render("TileLayout Interactive"))

	fmt.Fprintf(&sb, "%s\n", headStyle.Render("Keys"))
	_, isLayout := m.selectedNode().(*tl.TileLayout)
	if isLayout {
		for _, k := range []struct{ key, desc string }{
			{"h", "add tile (horiz layout)"},
			{"v", "add tile (vert layout)"},
			{"↑ / ↓", "select node"},
			{"q", "quit"},
		} {
			fmt.Fprintf(&sb, " %-8s %s\n", keyStyle.Render(k.key), dimStyle.Render(k.desc))
		}
	} else {
		for _, k := range []struct{ key, desc string }{
			{"h", "split horizontal"},
			{"v", "split vertical"},
			{"d", "delete selected"},
			{"↑ / ↓", "select node"},
			{"q", "quit"},
		} {
			fmt.Fprintf(&sb, " %-8s %s\n", keyStyle.Render(k.key), dimStyle.Render(k.desc))
		}
	}

	fmt.Fprintf(&sb, "\n%s\n", headStyle.Render("Preview Tree"))
	cp.renderTree(&sb, m.previewPanel, "")

	fmt.Fprintf(&sb, "\n%s\n",
		dimStyle.Render(fmt.Sprintf("%d tile(s)", len(m.tiles))),
	)

	return lipgloss.NewStyle().
		Width(cp.Size.Width - 1).
		Height(cp.Size.Height).
		BorderRight(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("237")).
		Padding(1, 1).
		Render(sb.String())
}

func (cp *ControlPanel) renderTree(sb *strings.Builder, l *tl.TileLayout, prefix string) {
	m := cp.model
	dirTag := "H"
	if l.Direction == tl.Vertical {
		dirTag = "V"
	}
	layoutStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("62"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)

	layoutLabel := dimStyle.Render("["+dirTag+"]") + " " + layoutStyle.Render(l.Name)
	if l.GetName() == m.selectedName {
		fmt.Fprintf(sb, "%s%s %s\n", prefix, selectedStyle.Render("▶"), layoutLabel)
	} else {
		fmt.Fprintf(sb, "%s  %s\n", prefix, layoutLabel)
	}
	for _, tile := range l.Tiles {
		if child, ok := tile.(*tl.TileLayout); ok {
			cp.renderTree(sb, child, prefix+"  ")
		} else {
			name := tile.GetName()
			if name == m.selectedName {
				fmt.Fprintf(sb, "%s    %s\n", prefix, selectedStyle.Render("▶ "+name))
			} else {
				fmt.Fprintf(sb, "%s    %s\n", prefix, dimStyle.Render("  "+name))
			}
		}
	}
}

// ─── Interactive model ────────────────────────────────────────────────────────

type InteractiveModel struct {
	root         *tl.TileLayout
	controlPanel *ControlPanel
	previewPanel *tl.TileLayout
	tiles        []*PreviewTile // flat ordered list of all leaf tiles
	selectedName string
	counter      int
	windowWidth  int
	windowHeight int
}

func NewInteractiveModel() *InteractiveModel {
	m := &InteractiveModel{}
	m.root = tl.NewRoot(tl.Horizontal)
	m.controlPanel = NewControlPanel(m)
	m.previewPanel = tl.NewTileLayout("Preview", tl.Vertical, tl.Size{Weight: 1.0})

	m.root.Add(m.controlPanel)
	m.root.Add(m.previewPanel)

	// Start with a single tile.
	m.counter++
	pt := NewPreviewTile(fmt.Sprintf("Tile-%d", m.counter))
	m.tiles = []*PreviewTile{pt}
	m.previewPanel.Add(pt)
	m.selectedName = pt.GetName()
	m.updateSelected()

	return m
}

// treeNodes returns all nodes (layouts + leaves) under the preview panel in pre-order.
func (m *InteractiveModel) treeNodes() []tl.Tile {
	var nodes []tl.Tile
	collectNodes(m.previewPanel, &nodes)
	return nodes
}

// selectedNode returns the node matching m.selectedName, or nil.
func (m *InteractiveModel) selectedNode() tl.Tile {
	for _, n := range m.treeNodes() {
		if n.GetName() == m.selectedName {
			return n
		}
	}
	return nil
}

// rebuildTileList rebuilds m.tiles by walking the preview tree in order.
func (m *InteractiveModel) rebuildTileList() {
	var newTiles []*PreviewTile
	var walk func(l *tl.TileLayout)
	walk = func(l *tl.TileLayout) {
		for _, tile := range l.Tiles {
			if child, ok := tile.(*tl.TileLayout); ok {
				walk(child)
			} else if pt, ok := tile.(*PreviewTile); ok {
				newTiles = append(newTiles, pt)
			}
		}
	}
	walk(m.previewPanel)
	m.tiles = newTiles
}

// selectNode moves selection by delta through treeNodes (clamped, no wrap).
func (m *InteractiveModel) selectNode(delta int) {
	nodes := m.treeNodes()
	if len(nodes) == 0 {
		return
	}
	idx := 0
	for i, n := range nodes {
		if n.GetName() == m.selectedName {
			idx = i
			break
		}
	}
	idx += delta
	if idx < 0 {
		idx = 0
	}
	if idx >= len(nodes) {
		idx = len(nodes) - 1
	}
	m.selectedName = nodes[idx].GetName()
	m.updateSelected()
}

// splitOrAdd handles h/v keys:
//   - layout selected: set layout direction to dir and add a new tile to it
//   - tile selected: wrap it in a new sub-layout (split)
func (m *InteractiveModel) splitOrAdd(dir tl.Direction) {
	node := m.selectedNode()
	if node == nil {
		return
	}

	if layout, ok := node.(*tl.TileLayout); ok {
		layout.Direction = dir
		m.counter++
		newTile := NewPreviewTile(fmt.Sprintf("Tile-%d", m.counter))
		newTile.Size.Weight = 1.0
		layout.Add(newTile)
		normalizeLayoutWeights(layout)
		m.rebuildTileList()
		m.selectedName = newTile.GetName()
		m.updateSelected()
		return
	}

	selected, ok := node.(*PreviewTile)
	if !ok {
		return
	}
	parent, ok := selected.GetParent().(*tl.TileLayout)
	if !ok || parent == nil {
		return
	}

	// Sub-layout inherits the weight of the tile it replaces in the parent.
	m.counter++
	sub := tl.NewTileLayout(
		fmt.Sprintf("Sub-%d", m.counter),
		dir,
		tl.Size{Weight: selected.GetSize().Weight},
	)
	parent.Replace(selected.GetName(), sub)

	s := selected.GetSize()
	s.Weight = 0.5
	selected.SetSize(s)
	sub.Add(selected)

	m.counter++
	newTile := NewPreviewTile(fmt.Sprintf("Tile-%d", m.counter))
	newTile.Size.Weight = 0.5
	sub.Add(newTile)

	m.rebuildTileList()
	m.selectedName = newTile.GetName()
	m.updateSelected()
}

// removeSelected deletes the currently selected tile (only leaf tiles can be removed).
func (m *InteractiveModel) removeSelected() {
	if len(m.tiles) <= 1 {
		return
	}
	node := m.selectedNode()
	selected, ok := node.(*PreviewTile)
	if !ok {
		return
	}
	parent, ok := selected.GetParent().(*tl.TileLayout)
	if !ok || parent == nil {
		return
	}

	// Pick the next selection before mutating the tree.
	nodes := m.treeNodes()
	selIdx := 0
	for i, n := range nodes {
		if n.GetName() == m.selectedName {
			selIdx = i
			break
		}
	}
	// Prefer the node before, fallback to the one after.
	nextIdx := selIdx - 1
	if nextIdx < 0 {
		nextIdx = selIdx + 1
	}

	parent.Remove(selected.GetName())

	// Collapse single-child sub-layouts.
	if parent != m.previewPanel && len(parent.Tiles) == 1 {
		grandparent, ok := parent.GetParent().(*tl.TileLayout)
		if ok && grandparent != nil {
			remaining := parent.Tiles[0]
			rs := remaining.GetSize()
			rs.Weight = parent.GetSize().Weight
			remaining.SetSize(rs)
			grandparent.Replace(parent.GetName(), remaining)
		}
	} else {
		normalizeLayoutWeights(parent)
	}

	m.rebuildTileList()

	newNodes := m.treeNodes()
	if nextIdx >= len(newNodes) {
		nextIdx = len(newNodes) - 1
	}
	if nextIdx >= 0 {
		m.selectedName = newNodes[nextIdx].GetName()
	} else if len(m.tiles) > 0 {
		m.selectedName = m.tiles[0].GetName()
	}
	m.updateSelected()
}

func (m *InteractiveModel) updateSelected() {
	for _, t := range m.tiles {
		t.SetFocused(t.GetName() == m.selectedName)
	}
}

func (m *InteractiveModel) relayout() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: m.windowWidth, Height: m.windowHeight}
	}
}

func (m *InteractiveModel) Init() tea.Cmd { return m.root.Init() }

func (m *InteractiveModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "h":
			m.splitOrAdd(tl.Horizontal)
			cmds = append(cmds, m.relayout())
		case "v":
			m.splitOrAdd(tl.Vertical)
			cmds = append(cmds, m.relayout())
		case "d":
			m.removeSelected()
			cmds = append(cmds, m.relayout())
		case "up", "k":
			m.selectNode(-1)
		case "down", "j":
			m.selectNode(1)
		}
	}

	_, cmd := m.root.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *InteractiveModel) View() string {
	return m.root.View()
}
