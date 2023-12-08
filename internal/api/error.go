package api

import "fmt"

type ErrorParentNodeDoesNotExist struct {
	ParentId string
}

func (e ErrorParentNodeDoesNotExist) Error() string {
	return fmt.Sprintf("parent node %s does not exist in the DoDAG, please register it first", e.ParentId)
}

type ErrorUnableToRegisterHandler struct {
	Pattern string
	err     error
}

func (e ErrorUnableToRegisterHandler) Error() string {
	return fmt.Sprintf("cannot register handler for pattern %s", e.Pattern)
}
