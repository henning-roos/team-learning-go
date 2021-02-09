package quiz

import (
	"io"
	"testing"

	"github.com/stretchr/testify/mock"
)

//var testQuestion = Question{
//	Question:     "Which language is this written in?",
//	RightAnswer:  "Go",
//	WrongAnswers: []string{"Python", "Java"},
//}

type QuizMock struct {
	mock.Mock
}

func (quizMock *QuizMock) ReadQuestionsFromJSON(jsonFile string) []Question {
	//TODO: understand why all calls gets the same parameter
	args := quizMock.Called(jsonFile)
	return args.Get(0).([]Question)
}

func (quizMock *QuizMock) GetAnswerMap(question Question) map[string]string {
	//TODO: Fix all the return values
	quizMock.Called()
	return nil
}

func (quizMock *QuizMock) GetUserInput(stdin io.Reader) (string, error) {
	quizMock.Called()
	return "", nil
}

func (quizMock *QuizMock) FormatQuestion(question Question, answerMap map[string]string) string {
	quizMock.Called()
	return ""
}

func (quizMock *QuizMock) Verify(question Question, answerMap map[string]string, userInput string) (bool, error) {
	quizMock.Called()
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

	quizMock.AssertCalled(t, "ReadQuestionsFromJSON")
	quizMock.AssertCalled(t, "GetAnswerMap")
	quizMock.AssertCalled(t, "GetUserInput")
	quizMock.AssertCalled(t, "FormatQuestion")
	quizMock.AssertCalled(t, "Verify")
}
