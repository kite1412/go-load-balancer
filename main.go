package main

import (
	"log"
	"net/http"

	"my.go/load-balancer/config"
	"my.go/load-balancer/instance"
)

var (
	logger = log.New(log.Writer(), "MAIN", log.Lshortfile)
)

func main() {
	conf, err := config.DefaultConfig("config.json", 8081, 8090)
	if err != nil {
		logger.Fatal(err)
	}
	roundRobinPool := instance.NewRoundRobinLoadBalancer(conf)
	http.ListenAndServe(":8080", roundRobinPool)
}
