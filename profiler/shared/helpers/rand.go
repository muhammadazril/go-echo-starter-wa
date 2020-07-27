package helpers

import "math/rand"

const letterSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterSet[rand.Int63()%int64(len(letterSet))]
	}
	return string(b)
}
