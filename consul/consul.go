package consul

import (
	"gitlab.com/gobang/bepkg/consul/grpc"
	"gitlab.com/gobang/bepkg/consul/http"
	"gitlab.com/gobang/bepkg/consul/models"
)

type agent struct {
}

func NewAgent() *agent {
	return &agent{}
}

func (a *agent) MustRegisterServiceWithGRPC(options *models.Options) {
	if options.HealthCheckInterval == 0 {
		options.HealthCheckInterval = 10
	}
	grpc.MustRegisterService(options)
}

func (a *agent) MustRegisterServiceWithHttp(options *models.Options) {
	if options.HealthCheckInterval == 0 {
		options.HealthCheckInterval = 10
	}
	http.MustRegisterService(options)
}
