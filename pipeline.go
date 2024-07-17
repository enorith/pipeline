package pipeline

import (
	"fmt"
	"sync"

	"github.com/alitto/pond"
	"github.com/enorith/pipeline/action"
)

const (
	WPDefaultSize   = 20
	WPDefaultBuffer = 1000
)

type Input struct {
	Type string         `json:"type"`
	From map[string]int `json:"from"`
}

type Node struct {
	Action      action.Action ``
	Inputs      []Input       `json:"inputs"`
	Outputs     []string      `json:"outputs"`
	Sigleton    bool          `json:"sigleton"`
	InvokeCount int           `json:"invokeCount"`
}

type Collection struct {
	nodes   map[string]*Node
	results map[string][]action.ActionParam
	mus     map[string]*sync.RWMutex
}

type PlayConfig struct {
	WPSize   int
	WPBuffer int
}

type PlayConfigFn func(config *PlayConfig)

// Play the pipeline
func (c *Collection) Play(config ...PlayConfigFn) {

	var conf = &PlayConfig{
		WPSize:   WPDefaultSize,
		WPBuffer: WPDefaultBuffer,
	}

	for _, fn := range config {
		fn(conf)
	}

	// 1. 查询目标节点（节点无输出，有引用输入）
	// 2. 执行目标节点，
	// 3. 查找引用的节点，执行后输出到目标节点
	// 4. 递归步骤 3
	pool := pond.New(conf.WPSize, conf.WPBuffer)
	var (
		targetNodeId string
		targetNode   *Node
	)

	for id, node := range c.nodes {
		if len(node.Outputs) == 0 && len(node.Inputs) > 0 {
			targetNodeId = id
			targetNode = node
			break
		}
	}

	defer pool.StopAndWait()

	c.callNode(targetNodeId, targetNode, pool)
}

func (c *Collection) callNode(id string, node *Node, pool *pond.WorkerPool) ([]action.ActionParam, error) {
	mu := c.mus[id]
	mu.Lock()
	defer mu.Unlock()

	if res, ok := c.results[id]; ok && node.Sigleton {
		return res, nil
	}

	var resChan = make(chan struct {
		outputs []action.ActionParam
		err     error
	})

	var params = make([]action.ActionParam, 0)
	if len(node.Inputs) > 0 {
		group := pool.Group()
		for _, input := range node.Inputs {
			for refNodeId, idx := range input.From {
				refNode := c.nodes[refNodeId]
				group.Submit(func() {
					ots, e := c.callNode(refNodeId, refNode, pool)
					resChan <- struct {
						outputs []action.ActionParam
						err     error
					}{
						outputs: ots,
						err:     e,
					}
				})

				result := <-resChan
				if result.err != nil {
					return nil, result.err
				}

				param := make(action.MargedParam, 0)

				if len(result.outputs) > idx {
					output := result.outputs[idx]
					if input.Type != output.GetType() {
						return nil, fmt.Errorf("invalid type of param[%d]: expected %s, %s given", idx, input.Type, output.GetType())
					}

					param = append(param, output)
				}

				if len(param) > 0 {
					params = append(params, param)
				}
			}
		}
		group.Wait()
	}
	node.InvokeCount++
	returns, e := node.Action.Handle(params...)
	if node.Sigleton {
		c.results[id] = returns
	}

	return returns, e
}

func NewCollection(nodes map[string]*Node) *Collection {
	mus := make(map[string]*sync.RWMutex)
	for id := range nodes {
		mus[id] = new(sync.RWMutex)
	}

	return &Collection{
		nodes:   nodes,
		results: make(map[string][]action.ActionParam),
		mus:     mus,
	}
}
