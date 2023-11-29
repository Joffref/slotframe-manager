package api

import "fmt"

type ErrorParentNodeDoesNotExist struct {
	ParentId string
}

func (e ErrorParentNodeDoesNotExist) Error() string {
	return fmt.Sprintf("parent node %s does not exist in the DoDAG, please register it first", e.ParentId)
}
