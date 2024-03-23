package consul_test

import (
	"testing"

	"github.com/armiariyan/bepkg/consul"
)

func TestGenerate(t *testing.T) {
	_ = consul.NewAgent()
}
