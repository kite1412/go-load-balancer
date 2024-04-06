// contains all load balancer implementations.
package loadbalancer

import (
	"log"

	"my.go/loadbalancer/config"
	"my.go/loadbalancer/instance"
)

var (
	logger = log.New(log.Writer(), "<LOADBALANCER>", log.Lshortfile)
)

// round robin load balancer implementation.
func RoundRobin() *instance.RoundRobinLoadBalancer {
	filepath, fErr := config.LBConfigAbs()
	ll, ul, pErr := config.LBLowUpLimitPort()
	
	if fErr != nil {
		logger.Fatal(fErr)
	}

	if pErr != nil {
		logger.Fatal(pErr)
	}

	conf, err := config.DefaultConfig(filepath, ll, ul)

	if err != nil {
		logger.Fatal(err)
	}

	lb := instance.NewRoundRobinLoadBalancer(conf)
	return lb
}