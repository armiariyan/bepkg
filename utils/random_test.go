package utils

import (
	"fmt"
	"testing"
)

func TestTimestampRandomInt(t *testing.T) {

	fmt.Println(TimestampRandomInt())
}

func TestRandomRefNumber(t *testing.T) {
	fmt.Println(GenerateRefNumber())
}
