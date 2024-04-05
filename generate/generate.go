package generate

import (
	"encoding/json"
	"fmt"
	"os"

	"my.go/load-balancer/config"
	"my.go/load-balancer/lberror"
)

func GeneratePort() (port int, err error) {
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
			return port, nil
		}
	}

	return -1, lberror.GeneratePortError("no available port.")
}

func contains(s []config.Backend, i int) bool {
	for _, b := range s {
		if b.Port == i {
			return true
		}
	}
	return false
}