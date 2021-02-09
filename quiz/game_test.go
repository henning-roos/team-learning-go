package quiz

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type QuizMock struct {
	mock.Mock
}

func (quizMock *QuizMock) ReadQuestionsFromJSON() {
	quizMock.Called()
}

func (quizMock *QuizMock) GetUserInput() {
	quizMock.Called()
}

func (quizMock *QuizMock) FormatQuestion() {
	quizMock.Called()
}

func (quizMock *QuizMock) GetAnswerMap() {
	quizMock.Called()
}

func (quizMock *QuizMock) Verify() {
	quizMock.Called()
}

func TestGame(t *testing.T) {
	quizMock := &QuizMock{}
	quizMock.On("ReadQuestionsFromJSON").Return(nil)
	quizMock.On("GetUserInput").Return(nil)
	quizMock.On("FormatQuestion").Return(nil)
	quizMock.On("GetAnswerMap").Return(nil)
	quizMock.On("Verify").Return(nil)

	quizMock.AssertCalled(t, "ReadQuestionsFromJSON")
	quizMock.AssertCalled(t, "GetUserInput")
	quizMock.AssertCalled(t, "FormatQuestion")
	quizMock.AssertCalled(t, "GetAnswerMap")
	quizMock.AssertCalled(t, "Verify")
}
