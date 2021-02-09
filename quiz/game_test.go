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
	quizMock.Called()
	return nil
}

func (quizMock *QuizMock) GetAnswerMap(question Question) map[string]string {
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
	quizMock.On("ReadQuestionsFromJSON").Return([]Question{testQuestion})
	quizMock.On("GetAnswerMap").Return(nil)
	quizMock.On("GetUserInput").Return("", nil)
	quizMock.On("FormatQuestion").Return("")
	quizMock.On("Verify").Return(false, nil)

	main(quizMock)

	quizMock.AssertCalled(t, "ReadQuestionsFromJSON")
	quizMock.AssertCalled(t, "GetAnswerMap")
	quizMock.AssertCalled(t, "GetUserInput")
	quizMock.AssertCalled(t, "FormatQuestion")
	quizMock.AssertCalled(t, "Verify")
}
