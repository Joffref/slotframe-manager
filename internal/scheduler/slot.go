package scheduler

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"log/slog"
)

type Slot struct {
	Number            int           // Number is the number of the slot
	EmittingChannels  []*graph.Node // Channels is a map of nodes indexed by their assigned channel in the slot
	ReceptionChannels []*graph.Node // Channels is a map of nodes indexed by their assigned channel in the slot
}

func NewSlot(number, numberOfChannels int) *Slot {
	emittingChannels := make([]*graph.Node, numberOfChannels)
	receptionChannels := make([]*graph.Node, numberOfChannels)
	for i := 0; i < numberOfChannels; i++ {
		emittingChannels[i] = nil
		receptionChannels[i] = nil
	}
	return &Slot{
		Number:            number,
		EmittingChannels:  emittingChannels,
		ReceptionChannels: receptionChannels,
	}
}

// AddNode adds a node to the slot
// ensuring that the node is not already in the slot
// and that the parent of the node is not emitting during this slot
// and that the parent of the node is not already listening during this slot
func (s *Slot) AddNode(node *graph.Node) bool {
	for _, n := range s.ReceptionChannels {
		if !isNotNeighbor(n, node) {
			return false
		}
	}

	for _, n := range s.EmittingChannels {
		if !isNotNeighbor(n, node) {
			return false
		}
	}

	for i, n := range s.EmittingChannels {
		if n == nil {
			slog.Debug(fmt.Sprintf("Adding node %s to slot %d", node.Id, s.Number))
			s.EmittingChannels[i] = node
			s.ReceptionChannels[i] = node.Parent
			node.Parent.ListeningSlots[s.Number] = i
			node.EmittingSlots[s.Number] = i
			return true
		}
	}
	return false
}

func isNotNeighbor(n *graph.Node, m *graph.Node) bool {
	if n == nil {
		return true
	}
	if n.Id == m.ParentId {
		return false
	}
	if n.Parent.Id == m.Id {
		return false
	}
	if n.Id == m.Parent.Id {
		return false
	}
	return true
}
