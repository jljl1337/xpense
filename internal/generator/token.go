package generator

import (
	"crypto/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewToken(length int) string {
	src := make([]byte, length)

	rand.Read(src)

	for i := range src {
		src[i] = charset[int(src[i])%len(charset)]
	}

	return string(src)
}
