package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type question struct {
	question     string
	rightAnswer  string
	wrongAnswers [2]string
}

type questionList struct {
    questions []question
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

// This function verifies that the answer is correct
func verify(testQuestion question, answer string) (bool, error) {

	if answer == testQuestion.rightAnswer {
		return true, nil
	}

	for _, value := range testQuestion.wrongAnswers {
		if answer == value {
			return false, nil
		}
	}

	return false, fmt.Errorf("The specified answer is invalid answer: %s", answer)
}

func readJson(jsonFile string) []question {
	file, _ := ioutil.ReadFile(jsonFile)

	data := []question

	_ = json.Unmarshal([]byte(file), &data)
	return nil
}
