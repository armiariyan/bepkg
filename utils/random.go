package utils

import (
	rand2 "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/oklog/ulid"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateThreadId() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	uniqueID := ulid.MustNew(ulid.Timestamp(t), entropy)
	return uniqueID.String()
}

//TimestampRandomInt create a simple string timestamp + 2 digit random integer
func TimestampRandomInt() string {

	return fmt.Sprintf("%d%d", time.Now().Unix(), rand.New(rand.NewSource(time.Now().UnixNano())).Intn(99))
}

//GenerateRefNumber create random ref number usually being used as Bank reference number
func GenerateRefNumber() string {
	randomByTime := time.Now().Format("20060102150405.9999Z07")
	result := strings.Replace(randomByTime, ".", "", -1)
	result = strings.Replace(result, "+", "", -1)
	result = strings.Replace(result, "Z", "", -1)

	return result + "0" + encodeToString(5)
}

func encodeToString(max int) string {
	b := make([]byte, max)
	_, _ = io.ReadAtLeast(rand2.Reader, b, max)

	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
