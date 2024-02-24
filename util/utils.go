package util

import "math/rand"

var alphabet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomString(n int) string {
	buf := make([]rune, n)

	for i := range buf {
		buf[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(buf)
}
