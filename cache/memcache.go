package cache

import (
	"time"

	"github.com/armiariyan/bepkg/logger"
	"github.com/bradfitz/gomemcache/memcache"
	"go.uber.org/zap"
)

type mcache struct {
	conn   *memcache.Client
	logger logger.Logger
}

// NewMemcache create new memcache client
func NewMemcache(servers []string) Keyval {
	mc := memcache.New(servers...)
	return &mcache{
		conn: mc,
	}
}

func (m *mcache) SetLogger(l logger.Logger) {
	m.logger = l
}

func (m *mcache) logInfo(method string, message interface{}) {
	if m.logger != nil {
		m.logger.Info("|",
			zap.String("_app_tag", "caching"),
			zap.String("_cache_type", "memcache"),
			zap.String("_cache_method", method),
			zap.Any("_message", message),
		)
	}
}

func (m *mcache) logError(method string, message interface{}) {
	if m.logger != nil {
		m.logger.Error("|",
			zap.String("_app_tag", "caching"),
			zap.String("_cache_type", "memcache"),
			zap.String("_cache_method", method),
			zap.Any("_message", message),
		)
	}
}

// Delete deletes the item with the provided key.
// Return nil byte if the item didn't already exist in the cache.
func (m *mcache) Get(key string) ([]byte, error) {

	item, err := m.conn.Get(key)

	if err == memcache.ErrCacheMiss {
		//Skip error if no value exist
		return nil, nil
	}

	if err != nil {
		m.logError("Get", err)
		return nil, err
	}

	m.logInfo("Get "+key, item)
	return item.Value, nil
}

// Add writes the given item, if no value already exists for its key.
// ErrNotStored is returned if that condition is not met.
func (m *mcache) Add(key string, val []byte, expiration time.Duration) (err error) {
	err = m.conn.Add(&memcache.Item{Key: key, Value: val})

	if err == nil {
		m.conn.Touch(key, int32(expiration*time.Second))
	}

	if err == memcache.ErrNotStored {
		//Skip error if value exist
		err = nil
	}

	if err != nil {
		m.logError("Add", err)
		return
	}

	m.logInfo("Add "+key, string(val))

	return
}

// Set writes the given item, unconditionally.
func (m *mcache) Set(key string, val []byte, expiration time.Duration) (err error) {
	err = m.conn.Set(&memcache.Item{Key: key, Value: val})
	if err == nil {
		m.conn.Touch(key, int32(expiration*time.Second))
	}
	if err != nil {
		m.logError("Set", err)
		return
	}
	m.logInfo("Set "+key, string(val))
	return
}

// Delete deletes the item with the provided key.
// return nil error if the item didn't already exist in the cache.
func (m *mcache) Delete(key string) error {
	err := m.conn.Delete(key)
	if err == memcache.ErrCacheMiss {
		//Skip error if no value exist
		return nil
	}
	if err != nil {
		m.logError("Get", err)
	}
	m.logInfo("Delete", key)
	return err
}
