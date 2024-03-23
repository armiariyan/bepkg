package postgres

import (
	"fmt"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var sqlConns = &SQLConn{DB: make(map[string]*SQLDB)}

// SQLConn is safe to use concurrently.
type SQLConn struct {
	DB  map[string]*SQLDB
	mux sync.Mutex
}

func (s *SQLConn) Get(id string) *SQLDB {
	s.mux.Lock()
	defer s.mux.Unlock()
	if conn, ok := s.DB[id]; ok {
		return conn
	}
	return nil
}

func (s *SQLConn) Set(id string, sqldb *SQLDB) {
	s.mux.Lock()
	s.DB[id] = sqldb
	s.mux.Unlock()
}

type Options struct {
	Username           string `json:"username"`
	Password           string `json:"password"`
	Name               string `json:"name"`
	Schema             string `json:"schema"`
	Host               string `json:"host"`
	Port               int    `json:"port"`
	MinIdleConnections int    `json:"minIdleConnections"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxLifetime        int    `json:"maxLifetime"`
	LogMode            bool   `json:"logmode"`
}

// Connect to mysql db based on id, and returns error if id is not exists / not valid.
// ex:
// id = master
// conf =
// map[string]string{
//   "username": "merchant",
//   "password": "merchant",
//   "name": "merchants",
//   "schema": "cico",
//   "host": "159.89.205.12",
//   "port": 5432,
//   "minIdleConnections": 10,
//   "maxOpenConnections": 30
// }
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

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		opt.Host,
		opt.Port,
		opt.Username,
		opt.Password,
		opt.Name,
		opt.Schema)
	//fmt.Println(dsn)

	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	conn.SetConnMaxLifetime(time.Duration(opt.MaxLifetime) * time.Second)
	conn.SetMaxOpenConns(opt.MaxOpenConnections)
	conn.SetMaxIdleConns(opt.MinIdleConnections)

	sqlConns.Set(id, &SQLDB{DB: conn, logMode: opt.LogMode})
	return sqlConns.Get(id), nil
}
