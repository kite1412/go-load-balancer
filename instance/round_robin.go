// contract and implementations of load balancer.
package instance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"my.go/loadbalancer/config"
	"my.go/loadbalancer/lberror"
)

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

// start http server with round robin load balancer attached
// and handle config update upon termination.
func (lb *RoundRobinLoadBalancer) Start(addr string) {
	con, err := lb.LbConfig.GetConfig()
	if err != nil {
		logger.Fatal(err)
		return
	}
	if !con.IsAlive {
		con.IsAlive = true
		lb.LbConfig.SetConfig(*con)
		
		go http.ListenAndServe(addr, lb)
		defer lb.terminate()

		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<- c
	}
}

func (lb *RoundRobinLoadBalancer) terminate() {
	con, err := lb.LbConfig.GetConfig()
	if err != nil {
		logger.Fatal(err)
		return
	}
	con.IsAlive = false
	lb.LbConfig.SetConfig(*con)
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

func contains(s []config.Backend, i int) bool {
	for _, b := range s {
		if b.Port == i {
			return true
		}
	}
	return false
}

// generate a port for server to be registered into
// load balancer's backend instances.
//
// url is typically the protocol followed by hostname
// e.g. http://localhost 
func (lb *RoundRobinLoadBalancer) GeneratePort(url string) (port int, err error) {
	cAbs, aErr := config.LBConfigAbs()
	if aErr != nil {
		logger.Println(aErr)
		return -1, aErr
	} 
	cont, rErr := os.ReadFile(cAbs)
	if rErr != nil {
		logger.Println(rErr)
		return -1, rErr
	}
	c := &config.Config{}
	if err := json.Unmarshal(cont, c); err != nil {
		logger.Println(err)
		return -1, err
	}
	ll, ul := c.PortLowerLimit, c.PortUpperLimit
	portSize := (ul - ll) + 1
	if len(c.Backends) == portSize {
		return -1, lberror.GeneratePortError("no available port.")
	}
	ports := make([]int, 0)
	for range portSize {
		ports = append(ports, ll)
		ll++
	}
	for _, port := range ports {
		if !contains(c.Backends, port) {
			new := config.Backend{
				Url: url + ":" + fmt.Sprint(port),
				Port: port,
			}
			lb.GetLbConfig().AddBackend(new)
			// send initial signal to newly added backend.
			go func(port int) {
				<- time.NewTicker(time.Second * 5).C
				logger.Println("<<<" + fmt.Sprint(port) + ":" + "INITIAL SIGNAL SENT>>>")
				lb.PingBackend(new)
			}(port)
			return port, nil
		}
	}

	return -1, lberror.GeneratePortError("no available port.")
}