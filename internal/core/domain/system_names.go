package domain

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numbers = []rune("0123456789")
)

func init() {
	// Seed once on startup
	rand.Seed(time.Now().UnixNano())
}

func GenerateName() string {
	// Format: LL-NNNN-L-NNNN (e.g., AB-1234-C-5678)
	l1 := pick(letters, 2)
	n1 := pick(numbers, 4)
	l2 := pick(letters, 1)
	n2 := pick(numbers, 4)

	return fmt.Sprintf("%s-%s-%s-%s", l1, n1, l2, n2)
}

func pick(source []rune, n int) string {
	result := make([]rune, n)
	for i := range result {
		result[i] = source[rand.Intn(len(source))]
	}
	return string(result)
}