package util

import (
	"math/rand"
	"strings"
	"time"
)

var characters = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomString(n int) string {
	var sb strings.Builder
	var charSize = len(characters)

	for i := 0; i < n; i++ {
		index := rand.Intn(charSize)
		sb.WriteByte(characters[index])
	}

	return sb.String()
}
