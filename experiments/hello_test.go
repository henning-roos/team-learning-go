package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testQuestion = Question{Question: "What is 1+1?", RightAnswer: "2", WrongAnswers: []string{"1", "54"}}

func TestHello(t *testing.T) {
	assert.Equal(t, hello(), "Hello World!", "Test hello world failed!!!")
}

func TestHelloFail(t *testing.T) {
	assert.False(t, !returnTrue(), "Test hello world negative test failed!!!")
}

func TestStruct(t *testing.T) {
	actualStruct := returnStruct()
	assert.Equal(t, actualStruct, testQuestion, "Didnt return expected")
	assert.Equal(t, actualStruct.Question, "What is 1+1?", "Not correct question")
}

func TestPrintHello(t *testing.T) {
	originalStdout := os.Stdout
	read, write, _ := os.Pipe()
	os.Stdout = write

	printHello()

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

func TestReadJSON(t *testing.T) {
	questions := readJSON("questions.json")
	assert.Equal(t, "What is 1+1?", questions[0].Question)
	assert.Equal(t, 3, len(questions))
}

func TestGetUserInput(t *testing.T) {
	var stdin bytes.Buffer
	var expected string = "2"
	stdin.Write([]byte(expected + "\n"))
	var input, _ = getUserInput(&stdin)
	assert.Equal(t, expected, input)
}

func TestGetUserInputError(t *testing.T) {
	var stdin bytes.Buffer
	stdin.Write([]byte("2"))
	var input, err = getUserInput(&stdin)
	assert.Error(t, err)
	assert.Equal(t, "", input)
}

func TestRandomizeAnswers(t *testing.T) {
	var answers = []string{"A", "B", "C"}
	var expected = []string{"B", "A", "C"}
	actual := randomizeAnswers(answers)
	assert.Equal(t, expected, actual)
}

func TestFormatQuestion(t *testing.T) {
	actual := formatQuestion(testQuestion)
	expected := "Question: What is 1+1?\n" +
		"1: 54\n" +
		"X: 1\n" +
		"2: 2\n" +
		"Answer: "
	assert.Equal(t, expected, actual)
}

// Question: What is 1+1?
// 1: 2
// X: 1
// 2: 53
// Your answer:

func TestMain(t *testing.T) {

}

type Foo struct {
	mock.Mock
}

func (myMock *Foo) Bar() {
	myMock.Called()
}

func TestFoo(t *testing.T) {
	myMock := &Foo{}
	myMock.On("Bar").Return(nil)

	myMock.Bar()
	myMock.AssertCalled(t, "Bar")
}