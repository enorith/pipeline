package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/alitto/pond"
	"github.com/enorith/pipeline/action"
)

const (
	WPDefaultSize          = 20
	WPDefaultBuffer        = 1000
	WPDefaultActionTimeout = 30 * time.Second
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

type NodeError struct {
	Err    error
	NodeId string
}

func (e NodeError) Error() string {
	return fmt.Sprintf("node %s error: %s", e.NodeId, e.Err.Error())
}

type Collection struct {
	nodes   map[string]*Node
	results map[string][]action.ActionParam
	mus     map[string]*sync.RWMutex
	mu      *sync.RWMutex
}

type PlayConfig struct {
	WPSize        int
	WPBuffer      int
	ActionTimeOut time.Duration
	TargetNodeId  string
}

type PlayConfigFn func(config *PlayConfig)

func PlayWithTargetId(id string) PlayConfigFn {
	return func(config *PlayConfig) {
		config.TargetNodeId = id
	}
}

func PlayWithActionTimeout(timeout time.Duration) PlayConfigFn {
	return func(config *PlayConfig) {
		config.ActionTimeOut = timeout
	}
}

// Play the pipeline
func (c *Collection) Play(config ...PlayConfigFn) ([]action.ActionParam, error) {

	var conf = &PlayConfig{
		WPSize:        WPDefaultSize,
		WPBuffer:      WPDefaultBuffer,
		ActionTimeOut: WPDefaultActionTimeout,
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
	if conf.TargetNodeId != "" {
		targetNodeId = conf.TargetNodeId
		targetNode = c.nodes[conf.TargetNodeId]
	} else {
		for id, node := range c.nodes {
			if len(node.Outputs) == 0 && len(node.Inputs) > 0 {
				targetNodeId = id
				targetNode = node
				break
			}
		}
	}

	defer pool.StopAndWait()

	return c.callNode(targetNodeId, targetNode, pool, conf)
}

type callResult struct {
	outputs             []action.ActionParam
	err                 error
	inputIdx, outputIdx int
	inputType, refNId   string
}

func (c *Collection) callNode(id string, node *Node, pool *pond.WorkerPool, conf *PlayConfig) ([]action.ActionParam, error) {
	mu := c.mu
	mu.RLock()
	if res, ok := c.results[id]; ok && node.Sigleton {
		mu.RUnlock()
		return res, nil
	}

	mu.RUnlock()

	inputLen := len(node.Inputs)
	var params = make([]action.ActionParam, inputLen)

	if inputLen > 0 {
		group := pool.Group()

		var resChan = make(chan callResult, inputLen)
		for i, input := range node.Inputs {
			for refNodeId, idx := range input.From {
				refNode := c.nodes[refNodeId]
				inputIdx := i
				inputType := input.Type
				outputIdx := idx
				refNId := refNodeId
				group.Submit(func() {
					ots, e := c.callNode(refNId, refNode, pool, conf)
					resChan <- callResult{
						outputs:   ots,
						err:       e,
						inputIdx:  inputIdx,
						inputType: inputType,
						outputIdx: outputIdx,
						refNId:    refNodeId,
					}
				})
			}
		}
		group.Wait()
		close(resChan)

		for result := range resChan {
			if result.err != nil {
				return nil, result.err
			}

			param := make(action.MargedParam, 0)
			if len(result.outputs) > result.outputIdx {
				output := result.outputs[result.outputIdx]
				if result.inputType != output.GetType() {
					return nil, fmt.Errorf("invalid type of param[%d]: expected %s, %s given", result.outputIdx, result.inputType, output.GetType())
				}

				param = append(param, output)
			}

			params[result.inputIdx] = param
		}
	}

	node.InvokeCount++
	ctx, cancel := context.WithTimeout(context.Background(), conf.ActionTimeOut)

	defer cancel()
	returns, e := node.Action.Handle(ctx, params...)

	if node.Sigleton && e == nil {
		mu.Lock()
		c.results[id] = returns
		mu.Unlock()
	}

	if e != nil {
		return nil, NodeError{
			NodeId: id,
			Err:    e,
		}
	}
	return returns, nil
}

func (c *Collection) GetNodes() map[string]*Node {
	return c.nodes
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
		mu:      new(sync.RWMutex),
	}
}
