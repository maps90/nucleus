package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringPrefix(t *testing.T) {
	random := RandomStringPrefix("test", 6)
	assert.NotEmpty(t, random)
}
