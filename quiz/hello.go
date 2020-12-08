package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Question struct {
	Question     string    `json:"question"`
	RightAnswer  string    `json:"rightAnswer"`
	WrongAnswers [2]string `json:"wrongAnswers"`
}

type questionList struct {
	questions []Question
}

func hello() string {

	return "Hello World!"
}

func returnTrue() bool {
	return true
}

func returnStruct() Question {
	return Question{Question: "What is 1+1?", RightAnswer: "2", WrongAnswers: [2]string{"1", "54"}}
}

func main() {
	fmt.Println(hello())
}

func countdown(out io.Writer) {
	fmt.Fprint(out, "3")
}

// This function verifies that the answer is correct
func verify(testQuestion Question, answer string) (bool, error) {

	if answer == testQuestion.RightAnswer {
		return true, nil
	}

	for _, value := range testQuestion.WrongAnswers {
		if answer == value {
			return false, nil
		}
	}

	return false, fmt.Errorf("The specified answer is invalid answer: %s", answer)
}

func readJson(jsonFile string) []Question {
	file, _ := ioutil.ReadFile(jsonFile)

	var data []Question

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func getUserInput(stdin io.Reader) string {
	reader := bufio.NewReader(stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "
	}
	return text, nil
}
