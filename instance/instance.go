// Contract and implementations of load balancer.
package instance

import (
	"log"
	"net/http"

	"my.go/loadbalancer/config"
)

var (
	logger = log.New(log.Writer(), "<INSTANCE>", log.Lshortfile)
)

type LoadBalancer interface {
	
	// Get list of available backends.
	GetBackends() ([]config.Backend, error)
	
	// Get next peer which decide where the request will be proceed to,
	// this is practically the implementation of existing load balancing algorithms.
	GetNextPeer() (config.Backend, error)
	
	// Get the config of the load balancer.
	GetLbConfig() config.LoadBalancerConfig
	
	// Send signal to a backend instance.
	PingBackend(config.Backend)
	
	ServeHTTP(http.ResponseWriter, *http.Request)
}
