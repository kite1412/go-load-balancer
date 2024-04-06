package main

import (
	"log"
	"net/http"

	"my.go/load-balancer/config"
	"my.go/load-balancer/instance"
	"my.go/load-balancer/state"
)

var (
	logger = log.New(log.Writer(), "MAIN", log.Lshortfile)
)

func main() {
	conf, err := config.DefaultConfig("config.json", 8081, 8090)
	if err != nil {
		logger.Fatal(err)
	}
	state.InitLoadBalancer(instance.NewRoundRobinLoadBalancer(conf))
	lb, _ := state.GetLoadBalancer()
	http.ListenAndServe(":8080", lb)
}
