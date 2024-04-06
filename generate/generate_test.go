package generate_test

import (
	"testing"
	"time"

	"my.go/load-balancer/config"
	"my.go/load-balancer/generate"
	"my.go/load-balancer/instance"
	"my.go/load-balancer/state"
)

func TestGeneratePort(t *testing.T) {
	conf, err := config.DefaultConfig("../config.json", 8081, 8090)
	if err != nil {
		t.Error(err)
	}
	state.InitLoadBalancer(instance.NewRoundRobinLoadBalancer(conf))
	port, err := generate.GeneratePort("http://localhost")
		if err != nil {
			t.Error(err)
		}
	t.Log(port)
	<-time.Tick(time.Second * 6)
}