package config_test

import (
	"my.go/loadbalancer/config"
	"testing"
)

const (
	filepath = "../config.json"
)

var (
	lbConfig config.LoadBalancerConfig
)

func TestInitDefault(t *testing.T) {
	defaultCon, err := config.DefaultConfig(filepath, 8081, 8090)
	if err != nil {
		t.Error(err)
	}
	lbConfig = defaultCon
	t.Log(lbConfig.GetConfig())
}

func TestGetConfig(t *testing.T) {
	_, err := lbConfig.GetConfig()
	if err != nil {
		t.Error(err)
	}
	t.Log("success getting config")
}

func TestAddBackends(t *testing.T) {
	lbConfig.AddBackend(config.Backend{
		Url: "http://localhost:8081",
		Port: 8081,
	})
	con, _ := lbConfig.GetConfig()
	t.Log(con.Backends)
}

func TestAddBulk(t *testing.T) {
	instance, ok := lbConfig.(*config.Default)
	if !ok {
		t.Error("Default not implements the the contract")
	}

	instance.AddBackends(
		config.Backend{
			Url: "http://localhost:8082",
			Port: 8082,
		},
		config.Backend{
			Url: "http://localhost:8083",
			Port: 8083,
		},
	)
}