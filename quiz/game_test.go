package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

type QuizMock struct {
	mock.Mock
}

func (quizMock *QuizMock) GetQuestions(configuration Configuration) []Question {
	args := quizMock.Called(configuration)
	return args.Get(0).([]Question)
}

func (quizMock *QuizMock) GetAnswerMap(question Question, randomizeAnswers bool) map[string]string {
	args := quizMock.Called(question, randomizeAnswers)
	return args.Get(0).(map[string]string)
}

func (quizMock *QuizMock) GetUserInput(stdin io.Reader) (string, error) {
	args := quizMock.Called(stdin)
	return args.String(0), args.Error(1)
}

func (quizMock *QuizMock) FormatQuestion(question Question, answerMap map[string]string) string {
	args := quizMock.Called(question, answerMap)
	return args.String(0)
}

func (quizMock *QuizMock) Verify(question Question, answerMap map[string]string, userInput string) (bool, error) {
	args := quizMock.Called(question, answerMap, userInput)
	return args.Bool(0), args.Error(1)
}

func (quizMock *QuizMock) FormatResult(numberCorrectAnswers int, numberQuestions int) string {
	args := quizMock.Called(numberCorrectAnswers, numberQuestions)
	return args.String(0)
}

func TestRun_AnswerCorrect(t *testing.T) {
	quizMock := &QuizMock{}
	quizMock.On("GetQuestions", mock.Anything).Return([]Question{testQuestion})
	quizMock.On("GetAnswerMap", mock.Anything, mock.Anything).Return(testAnswerMap)
	quizMock.On("GetUserInput", mock.Anything).Return("2", nil)
	quizMock.On("FormatQuestion", mock.Anything, mock.Anything).Return("Formatted question")
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	quizMock.On("FormatResult", mock.Anything, mock.Anything).Return("Formatted result")

	var stdin bytes.Buffer
	stdin.Write([]byte("dummy"))

	run(quizMock, &stdin)

	quizMock.AssertNumberOfCalls(t, "GetQuestions", 1)
	quizMock.AssertCalled(t, "GetAnswerMap", testQuestion, true)
	quizMock.AssertNumberOfCalls(t, "GetAnswerMap", 1)
	quizMock.AssertCalled(t, "GetUserInput", &stdin)
	quizMock.AssertNumberOfCalls(t, "GetUserInput", 1)
	quizMock.AssertCalled(t, "FormatQuestion", testQuestion, testAnswerMap)
	quizMock.AssertNumberOfCalls(t, "FormatQuestion", 1)
	quizMock.AssertCalled(t, "Verify", testQuestion, testAnswerMap, "2")
	quizMock.AssertNumberOfCalls(t, "Verify", 1)
	quizMock.AssertCalled(t, "FormatResult", 1, 1)
	quizMock.AssertNumberOfCalls(t, "FormatResult", 1)
}

func TestRun_AnswerWrong(t *testing.T) {
	quizMock := &QuizMock{}
	quizMock.On("GetQuestions", mock.Anything).Return([]Question{testQuestion})
	quizMock.On("GetAnswerMap", mock.Anything, mock.Anything).Return(testAnswerMap)
	quizMock.On("GetUserInput", mock.Anything).Return("1", nil)
	quizMock.On("FormatQuestion", mock.Anything, mock.Anything).Return("Formatted question")
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	quizMock.On("FormatResult", mock.Anything, mock.Anything).Return("Formatted result")

	var stdin bytes.Buffer
	stdin.Write([]byte("dummy"))

	run(quizMock, &stdin)

	quizMock.AssertNumberOfCalls(t, "GetQuestions", 1)
	quizMock.AssertCalled(t, "GetAnswerMap", testQuestion, true)
	quizMock.AssertNumberOfCalls(t, "GetAnswerMap", 1)
	quizMock.AssertCalled(t, "GetUserInput", &stdin)
	quizMock.AssertNumberOfCalls(t, "GetUserInput", 1)
	quizMock.AssertCalled(t, "FormatQuestion", testQuestion, testAnswerMap)
	quizMock.AssertNumberOfCalls(t, "FormatQuestion", 1)
	quizMock.AssertCalled(t, "Verify", testQuestion, testAnswerMap, "1")
	quizMock.AssertNumberOfCalls(t, "Verify", 1)
	quizMock.AssertCalled(t, "FormatResult", 0, 1)
	quizMock.AssertNumberOfCalls(t, "FormatResult", 1)
}

func TestRun_AnswerInvalid(t *testing.T) {
	quizMock := &QuizMock{}
	quizMock.On("GetQuestions", mock.Anything).Return([]Question{testQuestion})
	quizMock.On("GetAnswerMap", mock.Anything, mock.Anything).Return(testAnswerMap)
	quizMock.On("FormatQuestion", mock.Anything, mock.Anything).Return("Formatted question")
	quizMock.On("GetUserInput", mock.Anything).Return("bad input", nil).Once()
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(false, fmt.Errorf("mock error")).Once()
	quizMock.On("GetUserInput", mock.Anything).Return("2", nil).Once()
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
	quizMock.On("FormatResult", mock.Anything, mock.Anything).Return("Formatted result")

	var stdin bytes.Buffer
	stdin.Write([]byte("dummy"))

	run(quizMock, &stdin)

	quizMock.AssertNumberOfCalls(t, "GetQuestions", 1)
	quizMock.AssertCalled(t, "GetAnswerMap", testQuestion, true)
	quizMock.AssertNumberOfCalls(t, "GetAnswerMap", 1)
	quizMock.AssertCalled(t, "GetUserInput", &stdin)
	quizMock.AssertNumberOfCalls(t, "GetUserInput", 2)
	quizMock.AssertCalled(t, "FormatQuestion", testQuestion, testAnswerMap)
	quizMock.AssertNumberOfCalls(t, "FormatQuestion", 1)
	quizMock.AssertCalled(t, "Verify", testQuestion, testAnswerMap, "bad input")
	quizMock.AssertCalled(t, "Verify", testQuestion, testAnswerMap, "2")
	quizMock.AssertNumberOfCalls(t, "Verify", 2)
	quizMock.AssertCalled(t, "FormatResult", 1, 1)
	quizMock.AssertNumberOfCalls(t, "FormatResult", 1)
}

func TestRun_MultipleQuestions(t *testing.T) {
	quizMock := &QuizMock{}
	quizMock.On("GetQuestions", mock.Anything).Return([]Question{testQuestion, testQuestion2})
	quizMock.On("GetAnswerMap", mock.Anything, mock.Anything).Return(testAnswerMap).Once()
	quizMock.On("FormatQuestion", mock.Anything, mock.Anything).Return("Formatted question")
	quizMock.On("GetUserInput", mock.Anything).Return("2", nil).Once()
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
	quizMock.On("GetAnswerMap", mock.Anything, mock.Anything).Return(testAnswerMap2).Once()
	quizMock.On("GetUserInput", mock.Anything).Return("X", nil).Once()
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
	quizMock.On("FormatResult", mock.Anything, mock.Anything).Return("Formatted result")

	var stdin bytes.Buffer
	stdin.Write([]byte("dummy"))

	run(quizMock, &stdin)

	quizMock.AssertNumberOfCalls(t, "GetQuestions", 1)
	quizMock.AssertCalled(t, "GetAnswerMap", testQuestion, true)
	quizMock.AssertCalled(t, "GetAnswerMap", testQuestion2, true)
	quizMock.AssertNumberOfCalls(t, "GetAnswerMap", 2)
	quizMock.AssertCalled(t, "GetUserInput", &stdin)
	quizMock.AssertNumberOfCalls(t, "GetUserInput", 2)
	quizMock.AssertCalled(t, "FormatQuestion", testQuestion, testAnswerMap)
	quizMock.AssertCalled(t, "FormatQuestion", testQuestion2, testAnswerMap2)
	quizMock.AssertNumberOfCalls(t, "FormatQuestion", 2)
	quizMock.AssertCalled(t, "Verify", testQuestion, testAnswerMap, "2")
	quizMock.AssertCalled(t, "Verify", testQuestion2, testAnswerMap2, "X")
	quizMock.AssertNumberOfCalls(t, "Verify", 2)
	quizMock.AssertCalled(t, "FormatResult", 1, 2)
	quizMock.AssertNumberOfCalls(t, "FormatResult", 1)
}
