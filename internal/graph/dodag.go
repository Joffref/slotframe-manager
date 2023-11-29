package graph

type DoDAG struct {
	Root  *Node
	Nodes map[string]*Node // Nodes is a map of nodes indexed by their id
}

func NewDoDAG(root *Node) *DoDAG {
	return &DoDAG{
		Root:  root,
		Nodes: make(map[string]*Node),
	}
}

func (d *DoDAG) AddNode(node *Node) {
	node.Parent = d.Nodes[node.ParentId]
	d.Nodes[node.ParentId].Children = append(d.Nodes[node.ParentId].Children, node)
	node.NumberOfSlotNeeded = node.ETX
	d.Nodes[node.Id] = node
}

func (d *DoDAG) IsNode(id string) bool {
	_, ok := d.Nodes[id]
	return ok
}
