package instance_test

import (
	"load-balancer/config"
	"load-balancer/instance"
	"testing"
)

const (
	filepath = "../config.json"
)

var (
	roundRobin instance.LoadBalancer
)

func TestInitLB(t *testing.T) {
	conf, err := config.DefaultConfig(filepath, 100000, 1231231)
	if err != nil {
		t.Error(err)
	}
	roundRobin = instance.NewRoundRobinLoadBalancer(conf)
}

func TestGetNextPeer(t *testing.T) {
	for range 4 {
		b, err := roundRobin.GetNextPeer()
		if err != nil {
			t.Error(err)
		}
		t.Log(b)
	}
}