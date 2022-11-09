package selector

import (
	"errors"
	"github.com/dipper-iot/bozo/registry"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Next is a function that returns the next node
// based on the selector's strategy
type Next func() (*registry.Node, error)

// Random is a random strategy algorithm for node selection
func Random(services []*registry.Service) Next {
	nodes := make([]*registry.Node, 0, len(services))

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, errors.New("None Available")
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(services []*registry.Service) Next {
	nodes := make([]*registry.Node, 0, len(services))

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	var i = rand.Int()
	var mtx sync.Mutex

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, errors.New("None Available")
		}

		mtx.Lock()
		node := nodes[i%len(nodes)]
		i++
		mtx.Unlock()

		return node, nil
	}
}
