package util

import (
	"math/rand"
	"time"
)

func init() {
	// Initialize the random number generator
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}
func RandomCurrency() string {
	currencies := []string{USD, EUR, VND}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return RandomString(6) + "@" + RandomString(5) + ".com"
}
