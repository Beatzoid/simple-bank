package util

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomInt returns a random number between min (inclusive) and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString returns a random string of length n
func RandomString(stringLength int) string {
	var stringBuilder strings.Builder
	alphabetLength := len(alphabet)

	for i := 0; i < stringLength; i++ {
		randomLetter := alphabet[rand.Intn(alphabetLength)]
		stringBuilder.WriteByte(randomLetter)
	}

	return stringBuilder.String()
}

// RandomOwner returns a random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney returns a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency returns a random currency code
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	currenciesLength := len(currencies)

	return currencies[rand.Intn(currenciesLength)]
}

// RandomEmail returns a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
