package mysql

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var sqlConns = &SQLConn{DB: make(map[string]*SQLDB)}

// SQLConn is safe to use concurrently.
type SQLConn struct {
	DB  map[string]*SQLDB
	mux sync.Mutex
}

//Get a connection based on id (master|slave)
func (s *SQLConn) Get(id string) *SQLDB {
	s.mux.Lock()
	defer s.mux.Unlock()
	if conn, ok := s.DB[id]; ok {
		return conn
	}
	return nil
}

//Set set connection
func (s *SQLConn) Set(id string, sqldb *SQLDB) {
	s.mux.Lock()
	s.DB[id] = sqldb
	s.mux.Unlock()
}

//Options for connection
type Options struct {
	DSN                string `json:"dsn"`
	MinIdleConnections int    `json:"minIdleConnections"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxLifetime        int    `json:"maxLifetime"`
	LogMode            bool   `json:"logmode"`
}

// Connect to mysql db based on id, and returns error if id is not exists / not valid.
func Connect(id string, opt *Options) (*SQLDB, error) {

	// if previously established, reuse and ping
	conn := sqlConns.Get(id)
	if conn != nil {
		/*if conn.Ping() != nil {
			return newConnection(id)
		}*/
		return conn, nil
	}

	return connect(id, opt)
}

// Close close all established db connections
func Close() (err error) {
	if sqlConns == nil {
		return nil
	}
	for _, c := range sqlConns.DB {
		err = c.Close()
	}
	return err
}

func connect(id string, opt *Options) (*SQLDB, error) {

	conn, err := sqlx.Connect("mysql", opt.DSN)
	if err != nil {
		return nil, err
	}

	if opt.MaxLifetime > 0 {
		conn.SetConnMaxLifetime(time.Duration(opt.MaxLifetime) * time.Second)
	}

	if opt.MaxOpenConnections > 0 {
		conn.SetMaxOpenConns(opt.MaxOpenConnections)
	}

	if opt.MinIdleConnections > 0 {
		conn.SetMaxIdleConns(opt.MinIdleConnections)
	}

	sqlConns.Set(id, &SQLDB{DB: conn, logMode: opt.LogMode})
	return sqlConns.Get(id), nil
}
