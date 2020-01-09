package utils

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// GenerateKey function returns a string of random alphabetical characters of
// the passed length, consisting of varying letter cases.
func GenerateKey(length int) string {
	var generatedKey bytes.Buffer

	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Int() % (len(ASCII) - 1)
		makeLower := NewBoolGen().RandBool()

		// could be a letter or a string digit
		randChar := string(ASCII[randIndex])
		if makeLower {
			char := randChar[0]
			if !unicode.IsDigit(rune(char)) {
				generatedKey.WriteString(strings.ToLower(randChar))
			} else {
				generatedKey.WriteString(string(char))
			}
			continue
		}
		generatedKey.WriteString(randChar)
	}

	return generatedKey.String()
}
