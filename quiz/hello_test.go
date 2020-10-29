package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testQuestion = question{question: "What is 1+1?", rightAnswer: "2", wrongAnswers: [2]string{"1", "54"}}

func TestHello(t *testing.T) {
	assert.Equal(t, hello(), "Hello World!", "Test hello world failed!!!")
}

func TestHelloFail(t *testing.T) {
	assert.False(t, !returnTrue(), "Test hello world negative test failed!!!")
}

func TestStruct(t *testing.T) {
	actualStruct := returnStruct()
	assert.Equal(t, actualStruct, testQuestion, "Didnt return expected")
	assert.Equal(t, actualStruct.question, "What is 1+1?", "Not correct question")
}

func TestMain(t *testing.T) {
	originalStdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write

	main()

	write.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = originalStdout

	assert.Equal(t, "Hello World!\n", string(out))
}

func TestCountdown(t *testing.T) {
	buffer := &bytes.Buffer{}

	countdown(buffer)

	got := buffer.String()
	want := "3"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestWrongAnswer(t *testing.T) {
	answer := "1"
	actual, _ := verify(testQuestion, answer)
	assert.False(t, actual)
}

func TestCorrectAnswer(t *testing.T) {
	answer := "2"
	actual, _ := verify(testQuestion, answer)
	assert.True(t, actual)
}

func TestInvalidAnswer(t *testing.T) {
	var answer string = ""
	_, err := verify(testQuestion, answer)

	assert.Error(t, err)
}
