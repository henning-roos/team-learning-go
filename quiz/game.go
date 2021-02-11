package quiz

import "bytes"

func main(quiz QuizInterface) {
	questions := quiz.ReadQuestionsFromJSON("questions.json")

	//TODO, make a loop
	quiz.GetAnswerMap(questions[0])

	var stdin bytes.Buffer
	quiz.GetUserInput(&stdin)
}
