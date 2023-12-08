package graph

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/api"
	"sync"
	"time"
)

type Node struct {
	mtx      sync.Mutex
	Parent   *Node
	LastSeen time.Time
	Id       string `json:"id"`
	ParentId string `json:"parentId"`
	ETX      int    `json:"etx"`
	api.Slots
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %s (parent: %s, etx: %d, emittingSlots: %v, listeningSlots: %v)", n.Id, n.ParentId, n.ETX, n.EmittingSlots, n.ListeningSlots)
}

func NewNode(parentId string, id string, etx int, emittingSlots, listeningSlots map[int]int) *Node {
	if emittingSlots == nil {
		emittingSlots = make(map[int]int)
	}
	if listeningSlots == nil {
		listeningSlots = make(map[int]int)
	}
	return &Node{
		LastSeen: time.Now(),
		Id:       id,
		ParentId: parentId,
		ETX:      etx,
		Slots: api.Slots{
			EmittingSlots:  emittingSlots,
			ListeningSlots: listeningSlots,
		},
	}
}
func (n *Node) LockNode() func() {
	n.mtx.Lock()
	return func() {
		n.mtx.Unlock()
	}
}

func (n *Node) Rank() int {
	if n.Parent == nil {
		return 0
	}
	return n.Parent.Rank() + 1
}
