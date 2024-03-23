package grpc

import (
	"fmt"
	"sync"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

func SetupResolver() {
	resolver.Register(SetupBuilder())
}

type consulBuilder struct {
}

type consulResolver struct {
	address              string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	lastIndex            uint64
}

func SetupBuilder() resolver.Builder {
	return &consulBuilder{}
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (cr resolver.Resolver, err error) {
	host, port, name, err := parseTarget(fmt.Sprintf("%s/%s", target.Authority, target.Endpoint))

	if err != nil {
		return
	}

	solver := consulResolver{
		address:              fmt.Sprintf("%s%s", host, port),
		name:                 name,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		lastIndex:            0,
	}

	solver.wg.Add(1)
	go solver.watcher()

	cr = &solver

	return

}

func (cr *consulResolver) watcher() {
	config := api.DefaultConfig()
	config.Address = cr.address
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("error create consul client: %v\n", err)
		return
	}

	for {
		services, metaInfo, err := client.Health().Service(cr.name, cr.name, true, &api.QueryOptions{WaitIndex: cr.lastIndex})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		} else {
			cr.lastIndex = metaInfo.LastIndex
		}

		var addressList []resolver.Address
		for _, service := range services {
			addr := fmt.Sprintf("%v:%v", service.Service.Address, service.Service.Port)
			addressList = append(addressList, resolver.Address{Addr: addr})
		}

		state := resolver.State{
			Addresses: addressList,
		}

		cr.cc.UpdateState(state)
	}

}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOption) {
}

func (cr *consulResolver) Close() {
}
