package pipeline_test

import (
	"testing"

	"github.com/enorith/pipeline"
	"github.com/enorith/pipeline/action"
)

func TestPlay(t *testing.T) {
	coll := pipeline.NewCollection(map[string]*pipeline.Node{
		"input": {
			Action: action.Func(func() string {
				return "Foo"
			}),
			Outputs:  []string{"string"},
			Sigleton: true,
		},
		"input2": {
			Action: action.Func(func() string {
				return "Bar"
			}),
			Outputs:  []string{"string"},
			Sigleton: true,
		},
		"target": {
			Action: action.Func(func(s, s2 string) {
				t.Logf("play result target: %s, %s\n", s, s2)
			}),
			Inputs: []pipeline.Input{
				{
					Type: "string",
					From: map[string]int{
						"input": 0,
					},
				},
				{
					Type: "string",
					From: map[string]int{
						"input2": 0,
					},
				},
			},
		},
	})

	coll.Play()
}
