// hold app's internal states to be accessed
// by the rest of application's components.
package state

import (
	"my.go/load-balancer/instance"
	"my.go/load-balancer/lberror"
)

var (
	loadBalancer instance.LoadBalancer
)

// init load balancer instance once.
func InitLoadBalancer(newLb instance.LoadBalancer) {
	if loadBalancer == nil {
		loadBalancer = newLb
	}
}

// return the app's load balancer instance.
func GetLoadBalancer() (instance.LoadBalancer, error) {
	if loadBalancer == nil {
		return nil, lberror.StateNotInstantiatedError("load balancer instance is nil.")
	}
	return loadBalancer, nil
}