package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestSet(t *testing.T) {

	x := NewMemcache([]string{"0.0.0.0:9002"})
	err := x.Add("test", []byte("ini lagi"), 1*time.Hour)
	fmt.Println("Set")
	fmt.Println(err)
}

func TestAdd(t *testing.T) {

	x := NewMemcache([]string{"0.0.0.0:9002"})
	err := x.Add("test", []byte("ini isi test"), 1*time.Hour)
	fmt.Println("Add")
	fmt.Println(err)
}

func TestGet(t *testing.T) {

	x := NewMemcache([]string{"0.0.0.0:9002"})
	b, err := x.Get("test")
	fmt.Println(string(b))
	fmt.Println(err)
}

func TestDelete(t *testing.T) {

	x := NewMemcache([]string{"0.0.0.0:9002"})
	err := x.Delete("test")
	fmt.Println(err)
}
