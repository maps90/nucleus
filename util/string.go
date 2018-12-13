package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	// chars random
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// numbers random
	numbers = "0123456789"
)

//RandomString to get random string
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

//RandomStringPrefix to get random string with prefix
func RandomStringPrefix(prefix string, strlen int) string {
	result := make([]byte, strlen)
	data := []byte(prefix + RandomString(strlen))
	for i := 0; i < strlen; i++ {
		result[i] = data[rand.Intn(len(data))]
	}
	return string(result)
}

//RandomNumber to get random number
func RandomNumber(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = numbers[rand.Intn(len(numbers))]
	}
	return string(result)
}

// TrimWhiteSpace : trim excessive white space
func TrimWhiteSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// ToSentence : convert slice string into sentence
// eg : this, that, and the other
func ToSentence(words []string, andOrOr string) (response string) {
	var lastSentence string
	l := len(words)
	if l > 1 {
		lastSentence = words[l-1]
		words = words[:l-1]
	}

	if l > 1 {
		response = fmt.Sprintf("%v %v %v", strings.Join(words, ", "), andOrOr, lastSentence)
	} else {
		response = fmt.Sprintf("%v", strings.Join(words, ", "))
	}

	return
}

