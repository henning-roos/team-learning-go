package main

import "fmt"

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
