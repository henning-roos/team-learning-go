package quiz

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testQuestion = Question{Question: "What is 1+1?", RightAnswer: "2", WrongAnswers: []string{"1", "54"}}
var quiz = Quiz{}

func TestWrongAnswer(t *testing.T) {
	answer := "1"
	actual, _ := quiz.verify(testQuestion, answer)
	assert.False(t, actual)
}

func TestCorrectAnswer(t *testing.T) {
	answer := "2"
	actual, _ := quiz.verify(testQuestion, answer)
	assert.True(t, actual)
}

func TestInvalidAnswer(t *testing.T) {
	var answer string = ""
	_, err := quiz.verify(testQuestion, answer)

	assert.Error(t, err)
}

func TestReadQuestionsFromJSON(t *testing.T) {
	questions := quiz.readQuestionsFromJSON("questions.json")
	assert.Equal(t, "What is 1+1?", questions[0].Question)
	assert.Equal(t, 3, len(questions))
}

func TestGetUserInput(t *testing.T) {
	var stdin bytes.Buffer
	var expected string = "2"
	stdin.Write([]byte(expected + "\n"))
	var input, _ = quiz.getUserInput(&stdin)
	assert.Equal(t, expected, input)
}

func TestGetUserInputError(t *testing.T) {
	var stdin bytes.Buffer
	stdin.Write([]byte("2"))
	var input, err = quiz.getUserInput(&stdin)
	assert.Error(t, err)
	assert.Equal(t, "", input)
}

func TestRandomizeAnswers(t *testing.T) {
	var answers = []string{"A", "B", "C"}
	var expected = []string{"B", "A", "C"}
	actual := quiz.randomizeAnswers(answers)
	assert.Equal(t, expected, actual)
}

func TestFormatQuestion(t *testing.T) {
	actual := quiz.formatQuestion(testQuestion)
	expected := "Question: What is 1+1?\n" +
		"1: 54\n" +
		"X: 1\n" +
		"2: 2\n" +
		"Answer: "
	assert.Equal(t, expected, actual)
}
