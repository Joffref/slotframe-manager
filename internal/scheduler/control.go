package scheduler

import (
	"fmt"
	"github.com/Joffref/slotframe-manager/internal/graph"
	"log/slog"
	"time"
)

// ControlLoop is a function that runs in a goroutine and periodically checks if nodes are still alive
// If a node is not alive, it is removed from the DoDAG
func (s *Scheduler) ControlLoop() {
	for {
		if s.dodag.Root == nil { // If the root is nil, the DoDAG is in reset process
			slog.Debug("DoDAG has no root, waiting for a new one")
			time.Sleep(s.cfg.KeepAliveInterval)
			continue
		}
		fn := s.dodag.Root.LockNode()
		if !s.isAlive(s.dodag.Root) { // If the root is not alive, the DoDAG is reset
			slog.Debug("Root is not alive, resetting DoDAG")
			s.dodag.Root = nil
			s.dodag.Nodes = make(map[string]*graph.Node)
			slog.Debug(fmt.Sprintf("DoDAG after reset: %s", s.dodag.String()))
			time.Sleep(s.cfg.KeepAliveInterval)
			continue
		}
		fn()
		for _, node := range s.dodag.Nodes {
			slog.Debug(fmt.Sprintf("Checking if %s is alive", node.Id))
			fn := node.LockNode()
			if node.Parent == nil && node != s.dodag.Root { // If the node is not the root and has no parent, it is removed
				slog.Debug(fmt.Sprintf("%s has no parent, removing it", node.Id))
				removeNode(s.dodag, node.Id)
			}
			if !s.isAlive(node) { // If the node is not alive, it is removed
				slog.Debug(fmt.Sprintf("%s is not alive, removing it", node.Id))
				removeNode(s.dodag, node.Id)
			}
			fn()
		}
		slog.Debug(fmt.Sprintf("DoDAG after control loop: %s", s.dodag.String()))
		time.Sleep(s.cfg.KeepAliveInterval)
	}
}

func removeNode(dodag *graph.DoDAG, id string) {
	for _, n := range dodag.Nodes {
		if n.ParentId == id {
			n.Parent = nil
		}
	}
	delete(dodag.Nodes, id)
}

func (s *Scheduler) isAlive(node *graph.Node) bool {
	return node.LastSeen.Add(s.cfg.KeepAliveTimeout).After(time.Now())
}
