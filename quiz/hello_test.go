package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHello(t *testing.T) {
	assert.Equal(t, hello(), "Hello World!", "Test hello world failed!!!")
}

func TestHelloFail(t *testing.T) {
	assert.False(t, !returnTrue(), "Test hello world negative test failed!!!")
}

func TestStruct(t *testing.T) {
	actualStruct := returnStruct()
	assert.Equal(t, actualStruct, question{question: "What is 1+1?", rightAnswer: "2", wrongAnswers: [2]string{"1", "54"}}, "Didnt return expected")
	assert.Equal(t, actualStruct.question, "What is 1+1?", "Not correct question")
}

// func TestMain(t *testing.T) {
// 	testify.mock
// }

func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	countdown(buffer)

	got := buffer.String()
	want := "3"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
