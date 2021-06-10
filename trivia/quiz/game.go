package quiz

import (
	"fmt"
	"io"
)

// https://golangbyexample.com/print-output-text-color-console/
const colorReset string = "\033[0m"
const colorRed string = "\033[31m"
const colorGreen string = "\033[32m"

const randomizeAnswers bool = true

func Run(quiz QuizInterface, stdin io.Reader) error {
	configuration, err := quiz.ReadConfigurationFromYAML("resources/config.yaml")
	if err != nil {
		return err
	}

	questions, err := quiz.GetQuestions(configuration)
	if err != nil {
		return err
	}

	correctAnswers := 0

	for index, question := range questions {
		answerMap := quiz.GetAnswerMap(question, randomizeAnswers)

		fmt.Printf("%d/%d", index+1, len(questions))
		fmt.Println(quiz.FormatQuestion(question, answerMap))

		for {
			userInput, inputError := quiz.GetUserInput(stdin)
			if inputError != nil {
				return inputError
			}
			isAnswerCorrect, verificationError := quiz.Verify(question, answerMap, userInput)

			if verificationError == nil {
				if isAnswerCorrect {
					fmt.Printf("Your answer is %scorrect%s.\n\n", colorGreen, colorReset)
					correctAnswers++
				} else {
					fmt.Printf("Your answer is %swrong%s. The correct answer is '%s'\n\n", colorRed, colorReset, question.RightAnswer)
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
