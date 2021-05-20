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

var configuration = Configuration{
	QuestionFile: "questions.json",
	TriviaURL:    "https://opentdb.com/api.php?amount=10&type=multiple",
}

func run(quiz QuizInterface, stdin io.Reader) error {
	questions := quiz.GetQuestions(configuration)
	correctAnswers := 0

	for _, question := range questions {
		answerMap := quiz.GetAnswerMap(question, randomizeAnswers)

		fmt.Println(quiz.FormatQuestion(question, answerMap))

		for {
			userInput, inputError := quiz.GetUserInput(stdin)
			if inputError != nil {
				return inputError
			}
			isAnswerCorrect, verificationError := quiz.Verify(question, answerMap, userInput)

			if verificationError == nil {
				if isAnswerCorrect {
					fmt.Println("Your answer is correct.\n")
					correctAnswers++
				} else {
					fmt.Println("Your answer is wrong.\n")
					//questions.Rightanswer
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
