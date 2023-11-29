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

func (s *Scheduler) Register(parentId string, id string, etx int, input api.RegisterReponse) (*api.RegisterReponse, error) {
	if !s.dodag.IsNode(parentId) {
		return nil, api.ErrorParentNodeDoesNotExist{ParentId: parentId}
	}
	if s.dodag.IsNode(id) {
		slots := s.dodag.Nodes[id].Slots
		s.dodag.Nodes[id].LastSeen = time.Now()
		return &api.RegisterReponse{
			Slots: slots,
		}, nil
	}
	node := graph.NewNode(parentId, id, etx, input.Slots)
	s.dodag.AddNode(node)
	slots := s.scheduleNode(node)
	return &api.RegisterReponse{
		Slots: slots,
	}, nil
}

func NewScheduler(dodag *graph.DoDAG, rankingFunc func([]*graph.Node) []*graph.Node, frameSize int) *Scheduler {
	return &Scheduler{
		RankingFunc: rankingFunc,
		Frame:       NewFrame(frameSize),
		dodag:       dodag,
	}
}

func (s *Scheduler) Schedule() *Frame {
	nodes := make([]*graph.Node, 0, len(s.dodag.Nodes))
	for _, node := range s.dodag.Nodes {
		nodes = append(nodes, node)
	}
	nodes = s.RankingFunc(nodes)
	for _, node := range nodes {
		s.scheduleNode(node)
	}
	return s.Frame
}

func (s *Scheduler) scheduleNode(node *graph.Node) []int {
	for _, slot := range s.Frame.Slots {
		if !slot.AddNode(node) {
			continue
		}
		panic("TODO")
	}
	return nil
}
