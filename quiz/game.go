package quiz

func main(quiz QuizInterface) {
	questions := quiz.ReadQuestionsFromJSON("questions.json")

	//TODO, make a loop
	quiz.GetAnswerMap(questions[0])
}
