package quiz

func main(quiz QuizInterface) {
	questions := quiz.ReadQuestionsFromJSON("questions.json")
}
