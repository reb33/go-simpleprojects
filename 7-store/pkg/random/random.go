package random

import "strings"

import "math/rand"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenSessionId(n int) string {
	var result strings.Builder
	for range n {
		result.WriteString(string(alphabet[rand.Intn(len(alphabet))]))
	}
	return result.String()
}
