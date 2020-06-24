package common

import (
	"math/rand"
	"strings"
	"time"
	"unsafe"
)

var randSource rand.Source

const (
	bits     = 6
	mask     = 1<<bits - 1
	maxIndex = 63 / bits

	characters = "abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ23456789"
)

func init() {
	randSource = rand.NewSource(time.Now().UnixNano())
}

// RandString
//   - param n: suggest at least 10 to avoid conflict
// referred to https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func RandString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, randSource.Int63(), maxIndex; i >= 0; {
		if remain == 0 {
			cache, remain = randSource.Int63(), maxIndex
		}
		if idx := int(cache & mask); idx < len(characters) {
			b[i] = characters[idx]
			i--
		}
		cache >>= bits
		remain--
	}

	return strings.ToLower(*(*string)(unsafe.Pointer(&b)))
}
