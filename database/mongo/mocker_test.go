package mongo

import (
	"errors"
	"fmt"
	"testing"
)

func TestMockInterface(t *testing.T) {
	m := &Mock{}

	var x Client
	x = m
	fmt.Println(x)
}

func TestMockFindOne(t *testing.T) {
	m := &Mock{}
	type mockTrainer struct {
		Name string
	}
	o := mockTrainer{
		Name: "john",
	}

	var mo mockTrainer

	m.Stub("FindOne", M{}, MockResult{
		OutputVal: o,
		Error:     errors.New("test error"),
	})

	err := m.FindOne(M{}, &mo)

	fmt.Println(err)
	fmt.Println(mo)
}

func TestMockFindAll(t *testing.T) {
	m := &Mock{}
	type mockTrainer struct {
		Name string
	}
	o := []mockTrainer{
		mockTrainer{
			Name: "john",
		},
		mockTrainer{
			Name: "doe",
		},
	}

	var mo []mockTrainer

	m.Stub("Find", M{}, MockResult{
		OutputVal: o,
		Error:     errors.New("test error"),
	})

	err := m.Find(M{}, &mo)

	fmt.Println(err)
	fmt.Println(mo)
}
