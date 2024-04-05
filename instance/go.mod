module my.go/load-balancer/instance

go 1.22.1

require (
    my.go/load-balancer/config v1.0.0
    my.go/load-balancer/lberror v1.0.0
)

replace (
    my.go/load-balancer/config => ../config
    my.go/load-balancer/lberror => ../lberror
)
