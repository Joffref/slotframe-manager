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

func (f *Frame) String() string {
	return fmt.Sprintf("Frame (version: %d, slots: %v)", f.Version, f.Slots)
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

func ComputeFrame(dodag *graph.DoDAG, frameSize int, chanSize int) (*Frame, error) {
	slog.Debug(fmt.Sprintf("Computing frame with %d slots and %d channels per slot", frameSize, chanSize))

	frame := NewFrame(frameSize, chanSize)
	nodes := rankNodes(dodag)

	for i := len(nodes); i > 0; i-- { // Reverse loop on nodes to start from the highest rank and go up in the graph
		slog.Debug(fmt.Sprintf("Giving slots to nodes of rank %d", i-1))
		for _, node := range nodes[i-1] {
			err := frame.GiveSlots(0, node)
			if err != nil {
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

	neededSlots := node.ETX

	slog.Debug(fmt.Sprintf("Giving slots to %s", node.Id))
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
				err := s.GiveSlots(i, node.Parent)
				if err != nil {
					return err
				}
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
