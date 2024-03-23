package cache

import (
	"fmt"
	"testing"
)

func TestCacheMock(t *testing.T) {
	m := &Mock{}
	m.StubGet = func() ([]byte, error) {
		return []byte(`lorem pssum`), nil
	}

	var obj Keyval
	obj = m

	b, err := obj.Get("test")
	fmt.Println(string(b))
	fmt.Println(err)
}
