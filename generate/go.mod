module my.go/load-balancer/generate

go 1.22.1

require (
    my.go/load-balancer/config v1.0.0
    my.go/load-balancer/lberror v1.0.0
    my.go/load-balancer/state v1.0.0
)

require (
    my.go/load-balancer/instance v1.0.0 // indirect
)

replace (
    my.go/load-balancer/config => ../config
    my.go/load-balancer/lberror => ../lberror
    my.go/load-balancer/state => ../state
    my.go/load-balancer/instance => ../instance
)