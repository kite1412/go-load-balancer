// contract and implementations of load balancer.
package instance

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"my.go/load-balancer/config"
)

var (
	logger = log.New(log.Writer(), "INSTANCE", log.Lshortfile)
)

type LoadBalancer interface {
	// get list of available backends.
	GetBackends() ([]config.Backend, error)
	// get next peer which decide where the request will be proceed to,
	// this is practically the implementation of existing load balancing algorithms.
	GetNextPeer() (config.Backend, error)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// load balancer with round robin implementation.
type RoundRobinLoadBalancer struct {
	Config              config.LoadBalancerConfig
	currentBackendIndex int
	m                   sync.RWMutex
}

func (lb *RoundRobinLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend, err := lb.GetNextPeer()
	if err != nil {
		w.Write(JsonResponse(http.StatusInternalServerError, "can't serve requests at the moment"))
		logger.Fatal(err)
	}

	backend.ReverseProxy.ServeHTTP(w, r)
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

	for i := range backends {
		if i == lb.currentBackendIndex {
			lb.m.Lock()
			lb.currentBackendIndex = (i + 1) % len(backends)
			lb.m.Unlock()
			return backends[i], nil
		}
	}

	return backends[0], nil
}

func parseUrl(raw string) *url.URL {
	parsed, err := url.Parse(raw)
	if err != nil {
		logger.Fatal(err)
	}
	return parsed
}

func NewRoundRobinLoadBalancer(config config.LoadBalancerConfig) *RoundRobinLoadBalancer {
	rrlb := &RoundRobinLoadBalancer{
		Config: config,
		m:      sync.RWMutex{},
	}
	c, err := rrlb.Config.GetConfig()
	if err == nil {
		for i := range c.Backends {
			this := &c.Backends[i]
			this.ReverseProxy = httputil.NewSingleHostReverseProxy(parseUrl(this.Url))
			c.Backends[i] = *this
		}
		rrlb.Config.SetConfig(*c)
	}
	return rrlb
}
