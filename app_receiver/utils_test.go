package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getString(t *testing.T) {
	got := getHelloWorldString()

	assert.Equal(t, "Hello world", got)
}
