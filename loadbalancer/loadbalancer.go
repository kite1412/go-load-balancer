// Contains all load balancer implementations.
package loadbalancer

import (
	"fmt"
	"log"

	"my.go/loadbalancer/config"
	"my.go/loadbalancer/instance"
	"my.go/loadbalancer/lberror"
)

var (
	logger = log.New(log.Writer(), "<LOADBALANCER>", log.Lshortfile)
)

// Round robin load balancer implementation with default settings.
//
// look config.LBConfigAbs and config.LBLowUpLimitPort for setting up
// the environment.
func DefaultRoundRobin() *instance.RoundRobinLoadBalancer {
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

// Get the default load balancer's config.
func getDefaultConfig() (*config.Config, error) {
	filepath, fErr := config.LBConfigAbs()

	if fErr != nil {
		logger.Fatal(fErr)
	}

	return config.Read(filepath)
}

// Set new config to the default load balancer's config.
func setDefaultConfig(cfg config.Config) error {
	filepath, fErr := config.LBConfigAbs()

	if fErr != nil {
		logger.Fatal(fErr)
	}

	return config.Write(filepath, cfg)
}

// Register backend into default load balancer's backend instances.
//
// the returned port should be used by registrator.
// url typically the protocol followed by hostname,
// e.g. http://localhost
func DefaultRegister(url string) (port int, err error) {
	con, err := getDefaultConfig()

	if err != nil {
		return -1, err
	}

	if !con.IsAlive {
		return -1, lberror.GeneratePortError("can't register now, load balancer server is not alive.")
	}

	ports := make([]int, 0)
	ll, ul := con.PortLowerLimit, con.PortUpperLimit
	portRange := (ul - ll) + 1

	for range portRange {
		ports = append(ports, ll)
		ll++
	}

	// reuse unactive registered backend first.
	for i, b := range con.Backends {
		if !b.IsAlive {
			con.Backends[i].Url = url + ":" + fmt.Sprint(b.Port)
			setDefaultConfig(*con)
			return b.Port, nil
		}
	}

	for _, p := range ports {
		if !instance.ContainPort(con.Backends, p) {
			new := config.Backend{
				Url:  url + ":" + fmt.Sprint(p),
				Port: p,
			}
			con.Backends = append(con.Backends, new)
			if err := setDefaultConfig(*con); err != nil {
				return -1, err
			}
			return p, nil
		}
	}

	return -1, lberror.GeneratePortError("all ports are occupied.")
}
