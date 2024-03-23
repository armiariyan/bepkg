package consul_test

import (
	"gitlab.com/gobang/bepkg/consul"
	"testing"
)

func TestGenerate(t *testing.T) {
	_ = consul.NewAgent()
}
