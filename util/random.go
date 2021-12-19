package util

import (
	"math/rand"
	"strings"
	"time"
)

var (
	randBytes []byte
	random    = rand.New(rand.NewSource(time.Now().Unix()))
)

func init() {
	for i := '0'; i <= '9'; i++ {
		randBytes = append(randBytes, byte(i))
	}
	for i := 'a'; i <= 'z'; i++ {
		randBytes = append(randBytes, byte(i))
	}
	for i := 'A'; i <= 'Z'; i++ {
		randBytes = append(randBytes, byte(i))
	}
}

// Setup randomBytes, default is [0-9,a-z,A-Z].
func SetRandomBytes(b []byte) {
	randBytes = make([]byte, len(b))
	copy(randBytes, b)
}

// Return n length random string in range of randBytes.
func RandomString(n int) string {
	var str strings.Builder
	for i := 0; i < n; i++ {
		str.WriteByte(randBytes[random.Intn(len(randBytes))])
	}
	return str.String()
}

// Return n length random number string in range of [0-9].
func RandomNumber(n int) string {
	var str strings.Builder
	for i := 0; i < n; i++ {
		str.WriteByte(randBytes[random.Intn(10)])
	}
	return str.String()
}
