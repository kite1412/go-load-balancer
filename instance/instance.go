// contract and implementations of load balancer.
package instance

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"my.go/load-balancer/config"
	"my.go/load-balancer/lberror"
)

var (
	logger = log.New(log.Writer(), "<INSTANCE>", log.Lshortfile)
)

type LoadBalancer interface {
	// get list of available backends.
	GetBackends() ([]config.Backend, error)
	// get next peer which decide where the request will be proceed to,
	// this is practically the implementation of existing load balancing algorithms.
	GetNextPeer() (config.Backend, error)
	// get the config of the load balancer.
	GetLbConfig() config.LoadBalancerConfig
	// send signal to a backend instance.
	PingBackend(config.Backend)
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// load balancer with round robin implementation.
type RoundRobinLoadBalancer struct {
	LbConfig              config.LoadBalancerConfig
	currentBackendIndex int
	m                   sync.RWMutex
}

func (lb *RoundRobinLoadBalancer) GetBackends() ([]config.Backend, error) {
	lb.m.RLock()
	defer lb.m.RUnlock()
	con, err := lb.LbConfig.GetConfig()
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
		return config.NilBackend(), lberror.NoBackendFoundError("no backend found.")
	}

	for i, b := range backends {
		if i == lb.currentBackendIndex {
			if b.IsAlive {
				lb.m.Lock()
				lb.currentBackendIndex = (i + 1) % len(backends)
				lb.m.Unlock()
				return backends[i], nil
			} 
			// delegate request to next server in case the selected one
			// is not alive.
			lb.currentBackendIndex = (i + 1) % len(backends)
		}
	}

	return config.NilBackend(), lberror.NoBackendFoundError("no server can handle the request at the moment.")
}

func (lb *RoundRobinLoadBalancer) GetLbConfig() config.LoadBalancerConfig {
	return lb.LbConfig
}

func (lb *RoundRobinLoadBalancer) PingBackend(backend config.Backend) {
	lbCon := lb.GetLbConfig()
	con, err := lbCon.GetConfig()
	if err != nil {
		logger.Println(err)
		return
	}
	for i, b := range con.Backends {
		if b == backend {
			if _, err := http.Get(b.Url); err == nil {
				b.IsAlive = true				
			} else {
				b.IsAlive = false
			}
			con.Backends[i] = b
			lb.GetLbConfig().SetConfig(*con)
			return
		}
	}
}

func (lb *RoundRobinLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend, err := lb.GetNextPeer()
	if err != nil {
		w.Write(
			JsonResponse(
				http.StatusInternalServerError, 
				err.Error(),
			))
		logger.Fatal(err)
	}

	backend.ReverseProxy.ServeHTTP(w, r)
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
		LbConfig: config,
		m:      sync.RWMutex{},
	}
	c, err := rrlb.LbConfig.GetConfig()
	if err == nil {
		for i := range c.Backends {
			this := &c.Backends[i]
			this.ReverseProxy = httputil.NewSingleHostReverseProxy(parseUrl(this.Url))
			c.Backends[i] = *this
		}
		rrlb.LbConfig.SetConfig(*c)
	}
	return rrlb
}
