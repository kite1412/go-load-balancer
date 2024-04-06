package state_test

import (
	"testing"

	"my.go/load-balancer/config"
	"my.go/load-balancer/instance"
	"my.go/load-balancer/state"
)

func TestInitLb(t *testing.T) {
	conf, err := config.DefaultConfig("../config.json", 8081, 8090)
	if err != nil {
		t.Error(err)
	}
	state.InitLoadBalancer(
		instance.NewRoundRobinLoadBalancer(conf),
	)
	lb := state.GetLoadBalancer().(*instance.RoundRobinLoadBalancer)

	if lb == nil {
		t.Error("lb instance is nil")
	} 
	backends, _ := lb.GetBackends()
	t.Log("lb instance:", backends)
}