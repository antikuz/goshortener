package main

import (
	"math/rand"
	"time"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

func generateURL() string {
	rand.Seed(time.Now().UnixMicro())
	result := ""
	for i := 5; i > 0; i-- {
		result += string(chars[rand.Intn(len(chars))])
	}

	return result
}

func main() {
	print(generateURL())
}
