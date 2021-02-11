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

func run(quiz QuizInterface, stdin io.Reader) error {
	questions := quiz.ReadQuestionsFromJSON("questions.json")

	//TODO, make a loop
	question := questions[0]
	answer := quiz.GetAnswerMap(question)

	fmt.Println(quiz.FormatQuestion(question, answer))

	userInput, inputError := quiz.GetUserInput(stdin)
	if inputError != nil {
		return inputError
	}
	quiz.Verify(question, answer, userInput)

	return nil
}
