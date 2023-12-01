package graph

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	mtx            sync.Mutex
	Children       []*Node
	Parent         *Node
	LastSeen       time.Time
	Id             string      `json:"id"`
	ParentId       string      `json:"parentId"`
	ETX            int         `json:"etx"`
	EmittingSlots  map[int]int `json:"emittingSlots"`
	ListeningSlots map[int]int `json:"listeningSlots"`
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %s (parent: %s, etx: %d, emittingSlots: %v, listeningSlots: %v)", n.Id, n.ParentId, n.ETX, n.EmittingSlots, n.ListeningSlots)
}

func NewNode(parentId string, id string, etx int, emittingSlots, listeninSlots map[int]int) *Node {
	if emittingSlots == nil {
		emittingSlots = make(map[int]int)
	}
	if listeninSlots == nil {
		listeninSlots = make(map[int]int)
	}
	return &Node{
		Children:       nil,
		Parent:         nil,
		LastSeen:       time.Now(),
		Id:             id,
		ParentId:       parentId,
		ETX:            etx,
		EmittingSlots:  emittingSlots,
		ListeningSlots: listeninSlots,
	}
}

func (n *Node) LockNode() func() {
	n.mtx.Lock()
	return func() {
		n.mtx.Unlock()
	}
}

func (n *Node) AddChild(node *Node) {
	n.Children = append(n.Children, node)
}

func (n *Node) RemoveChild(node *Node) {
	for i, child := range n.Children {
		if child == node {
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return
		}
	}
}

func (n *Node) IsChildren(node *Node) bool {
	if n.Parent == nil {
		return false
	}
	if n.Parent != node {
		return n.Parent.IsChildren(node)
	}
	return true
}

func (n *Node) Rank() int {
	if n.Parent == nil {
		return 0
	}
	return n.Parent.Rank() + 1
}
