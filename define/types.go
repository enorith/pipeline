package define

import (
	"sync"
)

type Option struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type Instancer func(opts ...Option) (interface{}, error)

var (
	typeRegister = make(map[string]Instancer)
	mu           = new(sync.RWMutex)
)

func RegisterType(name string, instancer Instancer) {
	mu.Lock()
	defer mu.Unlock()
	typeRegister[name] = instancer
}

func GetInstance(name string, opts ...Option) (interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()

	instancer, ok := typeRegister[name]
	if !ok {
		return nil, UndefinedTypeError(name)
	}
	return instancer(opts...)
}
