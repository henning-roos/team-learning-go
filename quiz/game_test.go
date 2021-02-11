package quiz

import (
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

type QuizMock struct {
	mock.Mock
}

func (quizMock *QuizMock) ReadQuestionsFromJSON(jsonFile string) []Question {
	args := quizMock.Called(jsonFile)
	return args.Get(0).([]Question)
}

func (quizMock *QuizMock) GetAnswerMap(question Question) map[string]string {
	args := quizMock.Called(question)
	return args.Get(0).(map[string]string)
}

func (quizMock *QuizMock) GetUserInput(stdin io.Reader) (string, error) {
	args := quizMock.Called()
	return args.String(0), args.Error(1)
}

func (quizMock *QuizMock) FormatQuestion(question Question, answerMap map[string]string) string {
	args := quizMock.Called()
	return ""
}

func (quizMock *QuizMock) Verify(question Question, answerMap map[string]string, userInput string) (bool, error) {
	args := quizMock.Called()
	return false, nil
}

func TestGame(t *testing.T) {
	quizMock := &QuizMock{}
	quizMock.On("ReadQuestionsFromJSON", mock.Anything).Return([]Question{testQuestion})

	//Return(nil, mock.AnythingOfType("error"))

	quizMock.On("GetAnswerMap", mock.Anything)
	quizMock.On("GetUserInput", mock.Anything)
	quizMock.On("FormatQuestion", mock.Anything, mock.Anything)
	quizMock.On("Verify", mock.Anything, mock.Anything, mock.Anything)

	main(quizMock)

	quizMock.AssertCalled(t, "ReadQuestionsFromJSON", "questions.json")
	quizMock.AssertCalled(t, "GetAnswerMap")
	quizMock.AssertCalled(t, "GetUserInput")
	quizMock.AssertCalled(t, "FormatQuestion")
	quizMock.AssertCalled(t, "Verify")
}
