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

	// TODO: make a loop for the questions.json
	question := questions[0]
	answer := quiz.GetAnswerMap(question)

	fmt.Println(quiz.FormatQuestion(question, answer))

	for {
		userInput, inputError := quiz.GetUserInput(stdin)
		if inputError != nil {
			return inputError
		}
		result, verificationError := quiz.Verify(question, answer, userInput)

		if verificationError == nil {
			fmt.Printf("Result is: %t\n", result)
			break
		}
		fmt.Println(verificationError)
	}

	return nil
}
