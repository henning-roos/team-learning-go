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

const randomizeAnswers bool = true

func run(quiz QuizInterface, stdin io.Reader) error {
	questions := quiz.ReadQuestionsFromJSON("questions.json")
	correctAnswers := 0

	for _, question := range questions {
		answer := quiz.GetAnswerMap(question, randomizeAnswers)

		fmt.Println(quiz.FormatQuestion(question, answer))
		for {
			userInput, inputError := quiz.GetUserInput(stdin)
			fmt.Println("what was the option?")
			fmt.Println(userInput)
			if inputError != nil {
				return inputError
			}
			isAnswerCorrect, verificationError := quiz.Verify(question, answer, userInput)

			if verificationError == nil {
				fmt.Printf("Result is: %t\n", isAnswerCorrect)
				if isAnswerCorrect {
					correctAnswers++
				}
				break
			} else {
				fmt.Println(verificationError)
			}

		}
	}

	formattedResult := quiz.FormatResult(correctAnswers, len(questions))
	fmt.Println(formattedResult)

	return nil
}
