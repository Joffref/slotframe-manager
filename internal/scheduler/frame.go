package scheduler

import "github.com/Joffref/slotframe-manager/internal/graph"

type Frame struct {
	Slots map[int]*Slot // Slots is a map of slots indexed by their id
}

type Slot struct {
	Channels map[int]*graph.Node // Channels is a map of nodes indexed by their assigned channel in the slot
}

func NewFrame(size int) *Frame {
	slots := make(map[int]*Slot, size)
	for i := 0; i < size; i++ {
		slots[i] = nil
	}
	return &Frame{
		Slots: slots,
	}
}

func NewSlot(numberOfChannels int) *Slot {
	channels := make(map[int]*graph.Node, numberOfChannels)
	for i := 0; i < numberOfChannels; i++ {
		channels[i] = nil
	}
	return &Slot{
		Channels: channels,
	}
}

func (s *Slot) AddNode(node *graph.Node) bool {
	for _, slotNode := range s.Channels {
		if slotNode == nil {
			continue
		}
		if node.IsChildren(slotNode) || slotNode.IsChildren(node) {
			return false
		}
	}
	for i, slotNode := range s.Channels {
		if slotNode == nil {
			s.Channels[i] = node
			return true
		}
	}
	return false
}
