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
	quizMock.On("Bar").Return(nil)

	quizMock.AssertCalled(t, "Bar")
}