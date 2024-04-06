package generate

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"my.go/load-balancer/config"
	"my.go/load-balancer/lberror"
	"my.go/load-balancer/state"
)

var (
	logger = log.New(log.Writer(), "<GENERATE>", log.Lshortfile)
)

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
func GeneratePort(url string) (port int, err error) {
	cAbs, aErr := config.LBConfigAbs()
	if aErr != nil {
		fmt.Println(aErr)
		return -1, aErr
	} 
	cont, rErr := os.ReadFile(cAbs)
	if rErr != nil {
		fmt.Println(rErr)
		return -1, rErr
	}
	c := &config.Config{}
	if err := json.Unmarshal(cont, c); err != nil {
		fmt.Println(err)
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
			lb, err := state.GetLoadBalancer()
			// update config if load balancer instance exists.
			if err == nil {
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
			}
			return port, nil
		}
	}

	return -1, lberror.GeneratePortError("no available port.")
}