package main

import (
	"my.go/loadbalancer"
)

func main() {
	lb := loadbalancer.DefaultRoundRobin()
	lb.Start(":8080")
}
