package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testQuestion = Question{
	Question:     "Which language is this written in?",
	RightAnswer:  "Go",
	WrongAnswers: [2]string{"Python", "Java"},
}

var testQuestion2 = Question{
	Question:     "What is blue and yellow together? (using watercolors)",
	RightAnswer:  "Green",
	WrongAnswers: [2]string{"Red", "Black"},
}

var testAnswerMap = map[string]string{
	"1": "Java",
	"X": "Python",
	"2": "Go",
}

var testAnswerMap2 = map[string]string{
	"1": "Black",
	"X": "Red",
	"2": "Green",
}

var quiz = Quiz{}

func TestReadQuestionsFromJSON(t *testing.T) {
	questions := quiz.ReadQuestionsFromJSON("questions.json")
	assert.Equal(t, "What is 1+1?", questions[0].Question)
	assert.Equal(t, 3, len(questions))
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
		"1: Java\n" +
		"X: Python\n" +
		"2: Go\n" +
		"Answer: "
	assert.Equal(t, expectedQandA, actualQandA)
}

func TestGetAnswerMap(t *testing.T) {
	var randomizeAnswers bool = false
	// var copyOfTestQuestion Question
	question := quiz.ReadQuestionsFromJSON("questions.json")[0]
	// question := testQuestion
	// copyOfTestQuestion = question
	// copyOfTestQuestion.WrongAnswers = make([2]string, len(question.WrongAnswers))
	// copy(copyOfTestQuestion.WrongAnswers, question.WrongAnswers)

	quiz.GetAnswerMap(question, randomizeAnswers)
	//assert.Equal(t, testAnswerMap, actual)
	// assert.Equal(t, question, copyOfTestQuestion)
}

func TestWrongAnswer(t *testing.T) {
	userInput := "1"
	actual, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.False(t, actual)
	assert.Equal(t, nil, err)
}

func TestCaseSensitivityAnswer(t *testing.T) {
	userInput := "x"
	actual, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.False(t, actual)
	assert.Equal(t, nil, err)
}

func TestCorrectAnswer(t *testing.T) {
	userInput := "2"
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
