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
