package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testQuestion = Question{
	Question:     "Which language is this written in?",
	RightAnswer:  "Go",
	WrongAnswers: [3]string{"Python", "Java", "Ruby"},
}

var testQuestion2 = Question{
	Question:     "What is blue and yellow together? (using watercolors)",
	RightAnswer:  "Green",
	WrongAnswers: [3]string{"Red", "Black", "Pink"},
}

var testAnswerMap = map[string]string{
	"1": "Ruby",
	"2": "Java",
	"3": "Python",
	"4": "Go",
}

var testAnswerMap2 = map[string]string{
	"1": "Pink",
	"2": "Black",
	"3": "Red",
	"4": "Green",
}

var quiz = Quiz{}

func TestReadQuestionsFromJSON(t *testing.T) {
	questions := quiz.ReadQuestionsFromJSON("questions.json")
	assert.Equal(t, "What is blue and yellow together? (using watercolors)", questions[0].Question)
	assert.Equal(t, "Green", questions[0].RightAnswer)
	assert.Equal(t, "Red", questions[0].WrongAnswers[0])
	assert.Equal(t, 3, len(questions))
}

func TestReadQuestionsFromURL(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(scenario.expectedRespStatus)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()

	questions := quiz.ReadQuestionsFromURL(testServer.URL)
}

func TestGetUserInput(t *testing.T) {
	var stdin bytes.Buffer
	var expected string = "2"
	stdin.Write([]byte(expected + "\n"))
	var input, _ = quiz.GetUserInput(&stdin)
	assert.Equal(t, expected, input)
}

func TestGetUserInputWindows(t *testing.T) {
	var stdin bytes.Buffer
	var expected string = "2"
	stdin.Write([]byte(expected + "\r\n"))
	var input, _ = quiz.GetUserInput(&stdin)
	assert.Equal(t, expected, input)
}

func TestGetUserInputError(t *testing.T) {
	var stdin bytes.Buffer
	stdin.Write([]byte("2"))
	var input, err = quiz.GetUserInput(&stdin)
	assert.Error(t, err)
	assert.Equal(t, "", input)
}

func TestRandomizeAnswers(t *testing.T) {
	var answers = []string{"A", "B", "C"}
	var expected = []string{"B", "A", "C"}
	var randomizeAnswers bool = false
	actual := quiz.randomizeAnswers(answers, randomizeAnswers)
	assert.Equal(t, expected, actual)
}

func TestFormatQuestion(t *testing.T) {
	actualQandA := quiz.FormatQuestion(testQuestion, testAnswerMap)
	expectedQandA := "Question: Which language is this written in?\n" +
		"1: Ruby\n" +
		"2: Java\n" +
		"3: Python\n" +
		"4: Go\n" +
		"Answer: "
	assert.Equal(t, expectedQandA, actualQandA)
}

func TestGetAnswerMap(t *testing.T) {
	var randomizeAnswers bool = false
	question := testQuestion
	actual := quiz.GetAnswerMap(question, randomizeAnswers)
	assert.Equal(t, testAnswerMap, actual)
}

func TestWrongAnswer(t *testing.T) {
	userInput := "1"
	actual, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.False(t, actual)
	assert.Equal(t, nil, err)
}

func TestCorrectAnswer(t *testing.T) {
	userInput := "4"
	actual, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.True(t, actual)
	assert.Equal(t, nil, err)
}

func TestInvalidAnswer(t *testing.T) {
	var userInput string = ""
	_, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.Error(t, err)
}

func TestFormatResult(t *testing.T) {
	numberCorrectAnswers := 1
	numberQuestions := 2
	formattedResult := quiz.FormatResult(numberCorrectAnswers, numberQuestions)
	expected := "You got 1 of 2 correct answers."

	assert.Equal(t, expected, formattedResult)
}
