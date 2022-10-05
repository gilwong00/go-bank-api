package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// return random int between zero and max
func GetRandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // 0 => max - min
}

// returns random string
func GetRandomString(length int) string {
	var sb strings.Builder
	strLength := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(strLength)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// get random owner
func GetRandomOwner() string {
	return GetRandomString(6)
}

func GetRandomBalance() int64 {
	return GetRandomInt(0, 1000)
}

func GetCurrencyType() string {
	currencies := []string{
		"USD",
		"EUR",
		"CAD",
	}

	l := len(currencies)
	return currencies[rand.Intn(l)]
}
