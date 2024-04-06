module my.go/loadbalancer/instance

go 1.22.1

require (
    my.go/loadbalancer/config v1.0.0
    my.go/loadbalancer/lberror v1.0.0
)

replace (
    my.go/loadbalancer/config => ../config
    my.go/loadbalancer/lberror => ../lberror
)
