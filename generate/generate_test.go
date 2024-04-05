package generate_test

import (
	"testing"
	"my.go/load-balancer/generate"
)

func TestGeneratePort(t *testing.T) {
	port, err := generate.GeneratePort()
	if err != nil {
		t.Error(err)
	}
	t.Log(port)
}