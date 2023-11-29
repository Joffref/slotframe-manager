package scheduler

import (
	"github.com/Joffref/slotframe-manager/internal/graph"
	"github.com/Joffref/slotframe-manager/internal/utils"
	"time"
)

func ControlLoop(dodag *graph.DoDAG) {
	// TODO: implement
	for {
		for _, node := range dodag.Nodes {
			fn := node.LockNode()
			defer fn()
			if !IsAlive(*node) {
				continue
			}
			delete(dodag.Nodes, node.Id) // node is considered dead and is removed from the DoDAG
		}
		if !IsAlive(*dodag.Root) {
			panic("root node is dead")
		}
		time.Sleep(utils.KeepAliveInterval)
	}
}

func IsAlive(node graph.Node) bool {
	if node.LastSeen.Add(utils.KeepAliveTimeout).After(time.Now()) {
		return true
	}
	return false
}
