package scheduler

import (
	"github.com/Joffref/slotframe-manager/internal/graph"
	"github.com/Joffref/slotframe-manager/internal/utils"
	"time"
)

// ControlLoop is a function that runs in a goroutine and periodically checks if nodes are still alive
// If a node is not alive, it is removed from the DoDAG
func ControlLoop(dodag *graph.DoDAG) {
	for {
		if dodag.Root == nil { // If the root is nil, the DoDAG is in reset process
			time.Sleep(utils.KeepAliveInterval)
			continue
		}
		fn := dodag.Root.LockNode()
		if !isAlive(dodag.Root) { // If the root is not alive, the DoDAG is reset
			dodag.Root = nil
			dodag.Nodes = make(map[string]*graph.Node)
			time.Sleep(utils.KeepAliveInterval)
			continue
		}
		fn()
		for _, node := range dodag.Nodes {
			fn := node.LockNode()
			if node.Parent == nil && node != dodag.Root { // If the node is not the root and has no parent, it is removed
				delete(dodag.Nodes, node.Id)
			}
			if !isAlive(node) { // If the node is not alive, it is removed
				node.Parent.RemoveChild(node)
				delete(dodag.Nodes, node.Id)
			}
			fn()
		}
		time.Sleep(utils.KeepAliveInterval)
	}
}

func isAlive(node *graph.Node) bool {
	return node.LastSeen.Add(utils.KeepAliveTimeout).After(time.Now())
}
