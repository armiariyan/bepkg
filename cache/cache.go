package cache

import (
	"time"

	"gitlab.com/gobang/bepkg/logger"
)

//Keyval key value interface
type Keyval interface {
	SetLogger(l logger.Logger)
	Add(key string, val []byte, expiration time.Duration) error
	Set(key string, val []byte, expiration time.Duration) error
	Delete(key string) error
	Get(key string) ([]byte, error)
}

//Topology Server topology usually for HA setup
type Topology int

const (
	//Standalone single server/instance, default setting is standalone
	Standalone Topology = iota
	//Cluster with sharding, usually being used for redis or memcached
	Cluster
	//Sentinel for redis HA
	Sentinel
)

//SentinelConfig sentinel configuration data
type SentinelConfig struct {
	//Master or primary name
	PrimaryName string
	//format list of sentinel addresses IP:PORT,
	Addrs []string
}

//Config set config for cache
type Config struct {
	Timeout  time.Duration
	AuthPass string
	Topology Topology
	Sentinel SentinelConfig

	//format list of IP:PORT, use this variable for standalone or clustered server
	//Standalone server will only use Servers[0]
	Servers []string

	//Connection Pool size
	PoolSize int
}

//Mock cache mockers
type Mock struct {
	StubGet    func() ([]byte, error)
	StubSet    func() error
	StubAdd    func() error
	StubDelete func() error
}

//SetLogger mocker
func (m *Mock) SetLogger(l logger.Logger) {}

//Add Mocker
func (m *Mock) Add(key string, val []byte, expiration time.Duration) error { return m.StubAdd() }

//Set mocker
func (m *Mock) Set(key string, val []byte, expiration time.Duration) error { return m.StubSet() }

//Delete mocker
func (m *Mock) Delete(key string) error { return m.StubDelete() }

//Get mocker
func (m *Mock) Get(key string) ([]byte, error) { return m.StubGet() }
