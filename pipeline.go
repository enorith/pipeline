package pipeline

import "github.com/enorith/pipeline/action"

type Input struct {
	Type string
	From map[string]int
}

type Node struct {
	Action      action.Action
	Inputs      []Input
	Outputs     []string
	Sigleton    bool
	InvokeCount int
}

type Collection struct {
	nodes   map[string]Node
	results map[string][]action.ActionParam
}
