package common

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomId(length int) string {
	// generating random string of 6 length
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
