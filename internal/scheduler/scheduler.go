//nolint:typecheck
package scheduler

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/api"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"log/slog"
	"time"
)

type Scheduler struct {
	cfg   *Config
	frame *Frame
	dodag *graph.DoDAG
}

func NewScheduler(dodag *graph.DoDAG, cfg *Config) *Scheduler {
	slog.Debug("Creating scheduler")
	return &Scheduler{
		cfg:   cfg,
		frame: NewFrame(cfg.FrameSize, cfg.NumberCh),
		dodag: dodag,
	}
}

func (s *Scheduler) Register(parentId string, id string, etx int, input api.Slots) (*api.Slots, error) {
	if parentId == "0" {
		parentId = ""
	}
	if parentId != "" && !s.dodag.IsNode(parentId) {
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
		var unlockFns []func()
		currentVersion := s.frame.Version
		for _, node := range s.dodag.Nodes { // Reset the slots of each node
			unlockFns = append(unlockFns, node.LockNode())
			node.ListeningSlots = make(map[int]int)
			node.EmittingSlots = make(map[int]int)
		}
		frame, err := ComputeFrame(s.dodag, s.cfg.FrameSize, s.cfg.NumberCh)
		if err != nil {
			slog.Error(fmt.Sprintf("cannot compute frame: %v", err))
			continue
		}
		for _, fn := range unlockFns {
			fn()
		}
		for i, slot := range frame.Slots {
			if !isEq(s.frame.Slots[i].EmittingChannels, slot.EmittingChannels) {
				frame.Version = currentVersion + 1
				s.frame = frame
				slog.Debug(fmt.Sprintf("New frame: %s", s.frame.String()))
				break
			}
			if !isEq(s.frame.Slots[i].ReceptionChannels, slot.ReceptionChannels) {
				frame.Version = currentVersion + 1
				s.frame = frame
				slog.Debug(fmt.Sprintf("New frame: %s", s.frame.String()))
				break
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func isEq(a, b []*graph.Node) bool {
	if len(a) != len(b) {
		return false
	}
	for i, n := range a {
		if n != b[i] {
			return false
		}
	}
	return true
}

func (s *Scheduler) Version() int {
	return s.frame.Version
}
