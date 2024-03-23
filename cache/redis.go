package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/armiariyan/bepkg/logger"
	"github.com/mediocregopher/radix/v3"
)

type rcache struct {
	client       radix.Client
	sentinelConn *radix.Sentinel
	logger       logger.Logger
}

// NewRedis create redis client
func NewRedis(cfg Config) (kv Keyval, err error) {
	var conn radix.Client
	var sentinelConn *radix.Sentinel
	var opts []radix.DialOpt
	var customConnFunc func(network, addr string) (radix.Conn, error)
	var TopologyType Topology
	var useOpts bool
	//Default Pool Size
	poolSize := 4

	// This is a ConnFunc which will set up a connection which is authenticated
	// and has a timeout on all operations
	if cfg.PoolSize != 0 {
		poolSize = cfg.PoolSize
		useOpts = true
	}
	if cfg.Timeout != 0 {
		opts = append(opts, radix.DialTimeout(cfg.Timeout))
		useOpts = true
	}
	if cfg.AuthPass != "" {
		opts = append(opts, radix.DialAuthPass(cfg.AuthPass))
		useOpts = true
	}
	customConnFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr, opts...)
	}

	TopologyType = cfg.Topology
	servers := cfg.Servers

	switch TopologyType {
	case Standalone:
		var popt []radix.PoolOpt
		if useOpts {
			popt = append(popt, radix.PoolConnFunc(customConnFunc))
		}
		// this pool will use our ConnFunc for all connections it creates
		conn, err = radix.NewPool("tcp", servers[0], poolSize, popt...)
	case Cluster:
		var copt []radix.ClusterOpt
		if useOpts {
			// this cluster will use the ClientFunc to create a pool to each node in the
			// cluster. The pools also use our customConnFunc, but have more connections
			copt = append(copt, radix.ClusterPoolFunc(func(network, addr string) (radix.Client, error) {
				return radix.NewPool(network, addr, poolSize, radix.PoolConnFunc(customConnFunc))
			}))
		}
		conn, err = radix.NewCluster(servers, copt...)
	case Sentinel:
		var sopt []radix.SentinelOpt
		if useOpts {
			sopt = append(sopt, radix.SentinelPoolFunc(func(network, addr string) (radix.Client, error) {
				return radix.NewPool(network, addr, poolSize, radix.PoolConnFunc(customConnFunc))
			}))
		}
		if len(cfg.Sentinel.Addrs) == 0 || cfg.Sentinel.PrimaryName == "" {
			return nil, errors.New("Incomplete redis sentinels configuration")
		}
		sentinelConn, err = radix.NewSentinel(cfg.Sentinel.PrimaryName, cfg.Sentinel.Addrs, sopt...)
		conn = sentinelConn
		// adr, _ := sentinelConn.Addrs()
		// fmt.Printf("Using Address: %s\n", adr)
	}

	if err != nil {
		return
	}

	kv = &rcache{
		client:       conn,
		sentinelConn: sentinelConn,
	}
	return
}

func (m *rcache) SetLogger(l logger.Logger) {
	m.logger = l
}

func (m *rcache) logInfo(message interface{}) {
	if m.logger != nil {
		m.logger.Info("redis-cache",
			logger.ToField("caller", logger.Caller(2)),
			logger.ToField("message", message),
		)
	}
}

func (m *rcache) logError(message interface{}) {
	if m.logger != nil {
		m.logger.Error("redis-cache",
			logger.ToField("caller", logger.Caller(2)),
			logger.ToField("message", message),
		)
	}
}

// Get the item with the provided key.
// Return nil byte if the item didn't already exist in the cache.
func (m *rcache) Get(key string) (rcv []byte, err error) {
	err = m.client.Do(radix.Cmd(&rcv, "GET", key))
	if err != nil {
		m.logError(fmt.Sprintf("%s %s %s", key, string(rcv), err.Error()))
		return
	}
	return
}

// Add writes the given item, if no value already exists for its key.
// ErrNotStored is returned if that condition is not met.
func (m *rcache) Add(key string, val []byte, expiration time.Duration) (err error) {

	args := []string{key, string(val)}

	//NX -- Only set the key if it does not already exist.
	args = append(args, "NX")

	if expiration != 0 {
		//EX seconds -- Set the specified expire time, in seconds.
		//PX milliseconds -- Set the specified expire time, in milliseconds.
		args = append(args, "EX", fmt.Sprintf("%d", int(expiration.Seconds())))
	}
	err = m.client.Do(radix.Cmd(nil, "SET", args...))
	if err != nil {
		m.logError(fmt.Sprintf("%s %s %s", key, string(val), err.Error()))
		return
	}

	return
}

// Set writes the given item, unconditionally.
func (m *rcache) Set(key string, val []byte, expiration time.Duration) (err error) {

	args := []string{key, string(val)}

	if expiration != 0 {
		//EX seconds -- Set the specified expire time, in seconds.
		//PX milliseconds -- Set the specified expire time, in milliseconds.
		args = append(args, "EX", fmt.Sprintf("%d", int(expiration.Seconds())))
	}

	err = m.client.Do(radix.Cmd(nil, "SET", args...))
	if err != nil {
		m.logError(fmt.Sprintf("%s %s %s", key, string(val), err.Error()))
		return
	}

	return
}

// Delete deletes the item with the provided key.
// return nil error if the item didn't already exist in the cache.
func (m *rcache) Delete(key string) (err error) {
	err = m.client.Do(radix.Cmd(nil, "DEL", key))
	if err != nil {
		m.logError(fmt.Sprintf("%s %s", key, err.Error()))
		return
	}
	return
}
