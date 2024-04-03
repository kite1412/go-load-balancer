// contract and implementations of load balancer.
package instance

import (
	"errors"
	"load-balancer/config"
	"sync"
)

type LoadBalancer interface {
	GetBackends() ([]config.Backend, error)
	GetNextPeer() (config.Backend, error)
	ProcessRequest(config.Backend)
}

// load balancer with round robin implementation.
type RoundRobinLoadBalancer struct {
	Config config.LoadBalancerConfig
	currentBackend config.Backend
	m sync.RWMutex
}

func NewRoundRobinLoadBalancer(config config.LoadBalancerConfig) *RoundRobinLoadBalancer {
	return &RoundRobinLoadBalancer{
		Config: config,
	}
}

func (lb *RoundRobinLoadBalancer) GetBackends() ([]config.Backend, error) {
	lb.m.RLock()
	defer lb.m.RUnlock()
	con, err := lb.Config.GetConfig()
	if err != nil {
		return nil, err
	}
	return con.Backends, nil
}

func (lb *RoundRobinLoadBalancer) GetNextPeer() (config.Backend, error) {
	backends, err := lb.GetBackends()
	if err != nil {
		return config.NilBackend(), err
	}

	if len(backends) == 0 {
		return config.NilBackend(), errors.New("no backend found")
	}

	for i, b := range backends {
		if b == lb.currentBackend {
			lb.m.Lock()
			defer lb.m.Unlock()
			if i + 1 < len(backends) {
				lb.currentBackend = backends[i + 1]
				return lb.currentBackend, nil
			} else {
				lb.currentBackend = backends[0]
				return lb.currentBackend, nil
			}
		}
	}

	lb.currentBackend = backends[0]
	return lb.currentBackend, nil
}

func (lb *RoundRobinLoadBalancer) ProcessRequest(config.Backend) {
	
}