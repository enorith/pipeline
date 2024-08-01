package pipeline_test

import (
	"testing"
	"time"

	"github.com/enorith/pipeline"
	"github.com/enorith/pipeline/action"
)

func TestPlay(t *testing.T) {
	coll := pipeline.NewCollection(map[string]*pipeline.Node{
		"input": {
			Action: action.Func(func() string {
				time.Sleep(time.Second)
				return "foo " + time.Now().Format("2006-01-02 15:04:05")
			}),
			Outputs:  []string{"string"},
			Sigleton: true,
		},
		"input2": {
			Action: action.Func(func() string {
				time.Sleep(time.Second)
				return "Bar " + time.Now().Format("2006-01-02 15:04:05")
			}),
			Outputs:  []string{"string"},
			Sigleton: true,
		},
		"input3": {
			Action: action.Func(func() string {
				time.Sleep(2 * time.Second)
				return "baz " + time.Now().Format("2006-01-02 15:04:05")
			}),
			Outputs:  []string{"string"},
			Sigleton: true,
		},
		"target": {
			Action: action.Func(func(s, s2, s3 string) {
				t.Logf("play result target: %s, %s, %s\n", s, s2, s3)
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
				{
					Type: "string",
					From: map[string]int{
						"input3": 0,
					},
				},
			},
		},
	})

	coll.Play()
	// js, _ := json.MarshalIndent(coll.GetNodes(), "", "  ")
	// fmt.Println(string(js))
}
