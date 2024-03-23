package grpc

import (
	"fmt"

	"github.com/armiariyan/bepkg/consul/models"
	"github.com/hashicorp/consul/api"
)

func MustRegisterService(cs *models.Options) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cs.ConsulAddress
	client, err := api.NewClient(consulConfig)
	if err != nil {
		panic(err)
	}

	agent := client.Agent()

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", cs.Name, cs.Address, cs.Port),
		Name:    cs.Name,
		Port:    cs.Port,
		Address: cs.Address,
		Tags:    []string{cs.Name},
		Check: &api.AgentServiceCheck{
			Interval: cs.HealthCheckInterval.String(),
			GRPC:     fmt.Sprintf("%v:%v/%v", cs.Address, cs.Port, cs.Name),
		},
	}

	fmt.Printf("Register to Consul gprc on %s\n", cs.ConsulAddress)
	if err := agent.ServiceRegister(reg); err != nil {
		panic(err)
	}

}
