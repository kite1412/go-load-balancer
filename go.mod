module my.go/loadbalancer/main

go 1.22.1

require my.go/loadbalancer v1.0.0

require (
	my.go/loadbalancer/config v1.0.0 // indirect
	my.go/loadbalancer/instance v1.0.0 // indirect
	my.go/loadbalancer/lberror v1.0.0 // indirect
)

replace (
	my.go/loadbalancer/config => ./config
	my.go/loadbalancer/instance => ./instance
	my.go/loadbalancer/lberror => ./lberror
	my.go/loadbalancer => ./loadbalancer
)