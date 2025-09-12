package generator

import (
	"crypto/rand"
)

func NewToken(length int, charset string) string {
	src := make([]byte, length)

	rand.Read(src)

	for i := range src {
		src[i] = charset[int(src[i])%len(charset)]
	}

	return string(src)
}
