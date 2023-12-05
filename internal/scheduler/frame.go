package scheduler

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"log/slog"
)

type Frame struct {
	Version int     // Version is the version of the frame
	Slots   []*Slot // Slots is a map of slots indexed by their id
}

type Slot struct {
	Number            int                 // Number is the number of the slot
	EmittingChannels  map[int]*graph.Node // Channels is a map of nodes indexed by their assigned channel in the slot
	ReceptionChannels map[int]*graph.Node // Channels is a map of nodes indexed by their assigned channel in the slot

}

func NewFrame(size, channelPerSlot int) *Frame {
	slots := make([]*Slot, size)
	for i := 0; i < size; i++ {
		slots[i] = NewSlot(i, channelPerSlot)
	}
	return &Frame{
		Slots: slots,
	}
}

func NewSlot(number, numberOfChannels int) *Slot {
	emittingChannels := make(map[int]*graph.Node, numberOfChannels)
	receptionChannels := make(map[int]*graph.Node, numberOfChannels)
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
		if n == node.Parent {
			return false
		}
	}
	for _, n := range s.EmittingChannels {
		if n == node.Parent {
			return false
		}
	}
	for i, n := range s.EmittingChannels {
		if n == nil {
			s.EmittingChannels[i] = node
			s.ReceptionChannels[i] = node.Parent
			node.Parent.ListeningSlots[s.Number] = i
			node.EmittingSlots[s.Number] = i
			return true
		}
	}
	return false
}

func ComputeFrame(dodag *graph.DoDAG, frameSize int, chanSize int) (*Frame, error) {
	slog.Debug(fmt.Sprintf("Computing frame with %d slots and %d channels per slot", frameSize, chanSize))
	frame := NewFrame(frameSize, chanSize)
	nodes := rankNodes(dodag)
	for i := len(nodes); i > 0; i-- { // Reverse loop on nodes to start from the highest rank and go up in the graph
		slog.Debug(fmt.Sprintf("Giving slots to nodes of rank %d", i-1))
		for _, node := range nodes[i-1] {
			if err := frame.GiveSlots(0, node); err != nil {
				return nil, err
			}
		}
	}
	return frame, nil
}

func (s *Frame) GiveSlots(startAt int, node *graph.Node) error {
	if node.Parent == nil {
		slog.Debug(fmt.Sprintf("%s has no parent, won't give it slots as it is the root", node.Id))
		return nil
	}
	slog.Debug(fmt.Sprintf("Giving slots to %s", node.Id))
	neededSlots := node.ETX
	slog.Debug(fmt.Sprintf("Node %s needs %d slots", node.Id, neededSlots))
	slog.Debug(fmt.Sprintf("Slots number: %d", len(s.Slots)))
	for i, slot := range s.Slots {
		if i < startAt {
			continue
		}
		if slot.AddNode(node) {
			slog.Debug(fmt.Sprintf("Giving slot %d to %s", i, node.Id))
			neededSlots--
		}
		if neededSlots == 0 {
			if node.Parent != nil {
				return s.GiveSlots(i, node.Parent)
			}
			return nil
		}
	}
	return fmt.Errorf("not enough slots")
}

func rankNodes(dodag *graph.DoDAG) map[int][]*graph.Node {
	nodes := make(map[int][]*graph.Node, 0) // nodes is a map of nodes indexed by their rank in the graph
	for _, node := range dodag.Nodes {
		nodes[node.Rank()] = append(nodes[node.Rank()], node)
	}
	return nodes
}
