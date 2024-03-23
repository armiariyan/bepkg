package cache

import (
	"fmt"
	"testing"
	"time"

	"gitlab.com/gobang/bepkg/logger"
)

var mockConf = Config{
	Topology: Sentinel,
	Sentinel: SentinelConfig{
		PrimaryName: "mymaster",
		Addrs: []string{
			"0.0.0.0:26000",
			"0.0.0.0:26001",
			"0.0.0.0:26002",
		},
	},
	PoolSize: 10,
}

var redistestConn, _ = NewRedis(mockConf)

func doTaskEvery(d time.Duration, f func(time.Time) bool) {
	for x := range time.Tick(d) {
		if ok := f(x); !ok {
			break
		}
	}
}

func TestRedisSet(t *testing.T) {

	x := redistestConn
	x.SetLogger(logger.New(logger.Options{
		Stdout: true,
	}))
	err := x.Add("test", []byte("ini lagi"), 1*time.Hour)
	fmt.Println("Set")
	fmt.Println(err)
}

func TestRedisAdd(t *testing.T) {
	x := redistestConn

	x.SetLogger(logger.New(logger.Options{
		Stdout: true,
	}))
	err := x.Add("test", []byte("ini isi test"), 1*time.Hour)

	fmt.Println("Add")
	fmt.Println(err)
}

func TestRedisGet(t *testing.T) {

	x := redistestConn
	x.SetLogger(logger.New(logger.Options{
		Stdout: true,
	}))
	b, err := x.Get("test")

	fmt.Println("Get")
	fmt.Println(string(b))
	fmt.Println(err)
}

func TestRedisDelete(t *testing.T) {

	x := redistestConn
	x.SetLogger(logger.New(logger.Options{
		Stdout: true,
	}))
	err := x.Delete("test")
	fmt.Println(err)
}

// func TestRedisTasked(t *testing.T) {
// 	err := redistestConn.Add("task", []byte("task one"), 1*time.Hour)
// 	if err != nil {
// 		t.Errorf("Err %v", err)
// 	}
// 	doTaskEvery(1*time.Second, func(tm time.Time) bool {
// 		fmt.Printf("Time: %v \n", tm)
// 		b, err := redistestConn.Get("task")
// 		if err != nil {
// 			t.Errorf("Err %v", err)
// 		}
// 		fmt.Println(string(b))

// 		return true
// 	})
// }
