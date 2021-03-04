package main

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
	correctAnswers := 0

	for _, question := range questions {
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
			fmt.Println("Result and Correct Answers")
			fmt.Println(result)
			fmt.Println(correctAnswers)

			fmt.Println(verificationError)
			if result == true {
				correctAnswers = correctAnswers + 1
			}
		}
	}

	result := quiz.FormatResult(correctAnswers, len(questions))
	fmt.Println(result)

	return nil
}
