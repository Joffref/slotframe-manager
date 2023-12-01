package graph

import (
	"fmt"
	"log/slog"
)

type DoDAG struct {
	Root  *Node
	Nodes map[string]*Node // Nodes is a map of nodes indexed by their id
}

func (d *DoDAG) String() string {
	return fmt.Sprintf("DoDAG (root: %s, nodes: %v)", d.Root, d.Nodes)
}

func NewDoDAG() *DoDAG {
	slog.Debug("Creating a new DoDAG")
	return &DoDAG{
		Nodes: make(map[string]*Node),
	}
}

func (d *DoDAG) AddNode(node *Node) {
	slog.Debug(fmt.Sprintf("Adding %s to the DoDAG", node.String()))
	if node.ParentId == "" {
		slog.Debug("Node has no parent, setting it as root")
		d.Root = node
		d.Nodes[node.Id] = node
		return
	}
	node.Parent = d.Nodes[node.ParentId]
	d.Nodes[node.ParentId].Children = append(d.Nodes[node.ParentId].Children, node)
	d.Nodes[node.Id] = node
	slog.Debug(fmt.Sprintf("Node %s added to the DoDAG", node.Id))
}

func (d *DoDAG) IsNode(id string) bool {
	_, ok := d.Nodes[id]
	return ok
}
