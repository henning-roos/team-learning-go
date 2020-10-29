package main

import (
	"fmt"
	"io"
)

type question struct {
	question     string
	rightAnswer  string
	wrongAnswers [2]string
}

func hello() string {

	return "Hello World!"
}

func returnTrue() bool {
	return true
}

func returnStruct() question {
	return question{question: "What is 1+1?", rightAnswer: "2", wrongAnswers: [2]string{"1", "54"}}
}

func main() {
	fmt.Println(hello())
}

func countdown(out io.Writer) {
	fmt.Fprint(out, "3")
}

func verify(testQuestion question, answer string) (bool, error) {

	if answer in testQuestion.rightAnswer or answer in testQuestion.wrongAnswer
	return answer == testQuestion.rightAnswer, nil
}
