package main

import (
	"my.go/loadbalancer/loadbalancer"
)

// create a new load balancer instance with round robin algorithm.
func main() {
	lb := loadbalancer.RoundRobin()
	lb.Start(":8080")
}