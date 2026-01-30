package gen

import (
	"fmt"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateStarName creates a name in the format: LL-NNNN-L-NNNN
func GenerateStarName() string {
	l1 := letters[r.Intn(len(letters))]
	l2 := letters[r.Intn(len(letters))]
	n1 := r.Intn(10000)
	l3 := letters[r.Intn(len(letters))]
	n2 := r.Intn(10000)

	return fmt.Sprintf("%c%c-%04d-%c-%04d", l1, l2, n1, l3, n2)
}
