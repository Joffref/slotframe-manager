package scheduler

import (
	"github.com/Joffref/slotframe-manager/internal/api"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"time"
)

type Scheduler struct {
	RankingFunc func([]*graph.Node) []*graph.Node
	Frame       *Frame
	dodag       *graph.DoDAG
}

func NewScheduler(dodag *graph.DoDAG, rankingFunc func([]*graph.Node) []*graph.Node, frameSize int) *Scheduler {
	return &Scheduler{
		RankingFunc: rankingFunc,
		Frame:       NewFrame(16, frameSize),
		dodag:       dodag,
	}
}

func (s *Scheduler) Register(parentId string, id string, etx int, input api.Slots) (*api.Slots, error) {
	if !s.dodag.IsNode(parentId) {
		return nil, api.ErrorParentNodeDoesNotExist{ParentId: parentId}
	}
	if s.dodag.IsNode(id) {
		s.dodag.Nodes[id].LastSeen = time.Now()
		return &api.Slots{
			EmittingSlots:  s.dodag.Nodes[id].EmittingSlots,
			ListeningSlots: s.dodag.Nodes[id].ListeningSlots,
		}, nil
	}
	node := graph.NewNode(parentId, id, etx, input.EmittingSlots, input.ListeningSlots)
	s.dodag.AddNode(node)
	return &api.Slots{
		EmittingSlots:  node.EmittingSlots,
		ListeningSlots: node.ListeningSlots,
	}, nil
}

func (s *Scheduler) Schedule() {
	for {
		currentVersion := s.Frame.Version
		s.Frame = ComputeFrame(s.dodag, 0, 0)
		s.Frame.Version = currentVersion + 1
	}
}

func (s *Scheduler) Version() int {
	return s.Frame.Version
}
