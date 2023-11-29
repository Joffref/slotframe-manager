package graph

import (
	"sync"
	"time"
)

type Node struct {
	mtx                sync.Mutex
	Children           []*Node
	Parent             *Node
	LastSeen           time.Time
	Id                 string `json:"id"`
	ParentId           string `json:"parentId"`
	ETX                int    `json:"etx"`
	NumberOfSlotNeeded int    `json:"-"`
	Slots              []int  `json:"slots"`
}

func NewNode(parentId string, id string, etx int, slots []int) *Node {
	return &Node{
		Children:           nil,
		Parent:             nil,
		LastSeen:           time.Now(),
		Id:                 id,
		ParentId:           parentId,
		ETX:                etx,
		NumberOfSlotNeeded: 0,
		Slots:              slots,
	}
}

func (n *Node) LockNode() func() {
	n.mtx.Lock()
	return func() {
		n.mtx.Unlock()
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

func (n *Node) AddNeededSlot() {
	if n.Parent == nil {
		return
	}
	n.Parent.NumberOfSlotNeeded += n.Parent.ETX
	n.Parent.AddNeededSlot()
}

func (n *Node) Rank() int {
	if n.Parent == nil {
		return 0
	}
	return n.Parent.Rank() + 1
}
