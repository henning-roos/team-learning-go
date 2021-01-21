package quiz

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testQuestion = Question{
	Question:     "Which language is this written in?",
	RightAnswer:  "Go",
	WrongAnswers: []string{"Python", "Java"},
}

var testAnswerMap = map[string]string{
	"1": "Java",
	"X": "Python",
	"2": "Go",
}

var quiz = Quiz{}

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
	actualQandA := quiz.formatQuestion(testQuestion, testAnswerMap)
	expectedQandA := "Question: Which language is this written in?\n" +
		"1: Java\n" +
		"X: Python\n" +
		"2: Go\n" +
		"Answer: "
	assert.Equal(t, expectedQandA, actualQandA)
}

func TestGetAnswerMap(t *testing.T) {
	actual := quiz.getAnswerMap(testQuestion)
	assert.Equal(t, testAnswerMap, actual)
}

func TestWrongAnswer(t *testing.T) {
	userInput := "x"
	actual, _ := quiz.verify(testQuestion, testAnswerMap, userInput)
	assert.False(t, actual)
}

func TestCorrectAnswer(t *testing.T) {
	userInput := "2"
	actual, _ := quiz.verify(testQuestion, testAnswerMap, userInput)
	assert.True(t, actual)
}

func TestInvalidAnswer(t *testing.T) {
	var userInput string = ""
	_, err := quiz.verify(testQuestion, testAnswerMap, userInput)

	assert.Error(t, err)
}

func TestPlayQuiz(t *testing.T) {
	questions := quiz.readQuestionsFromJSON("questions.json")
	//formattedQuestion = quiz.formatQuestion(questions[0])
	quiz.verify(questions[0], "54")
}
