package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		randCharsetIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset)-1)))
		b[i] = charset[randCharsetIndex.Int64()]
	}

	return string(b)
}
