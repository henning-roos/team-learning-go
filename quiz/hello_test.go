package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	assert.Equal(t, hello(), "Hello World!", "Test hello world failed!!!")
}
func TestHelloFail(t *testing.T) {
	assert.False(t, !returnTrue(), "Test hello world negative test failed!!!")
}
