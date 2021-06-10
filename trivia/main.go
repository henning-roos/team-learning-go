package main

import (
	"os"
	"trivia/quiz"
)

func main() {
	stdin := os.Stdin
	quizGame := &quiz.Quiz{}
	quiz.Run(quizGame, stdin)
}
