module my.go/load-balancer

go 1.22.1

require (
    my.go/load-balancer/instance v1.0.0
    my.go/load-balancer/config v1.0.0
)

replace (
    my.go/load-balancer/instance => ./instance
    my.go/load-balancer/config => ./config
)