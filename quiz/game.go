package quiz

import (
	"fmt"
	"io"
	"os"
)

func main() {
	stdin := os.Stdin
	quiz := &Quiz{}
	run(quiz, stdin)
}

func run(quiz QuizInterface, stdin io.Reader) {
	questions := quiz.ReadQuestionsFromJSON("questions.json")

	//TODO, make a loop
	question := questions[0]
	answer := quiz.GetAnswerMap(question)

	fmt.Println(quiz.FormatQuestion(question, answer))

	quiz.GetUserInput(stdin)
}
