module load-balancer/instance

go 1.22.1

require(
    load-balancer/config v1.0.0
)

replace(
    load-balancer/config => C:/go-project/go-load-balancer/config
)