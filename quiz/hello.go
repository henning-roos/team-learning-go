package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
)

type Question struct {
	Question     string   `json:"question"`
	RightAnswer  string   `json:"rightAnswer"`
	WrongAnswers []string `json:"wrongAnswers"`
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
	return Question{Question: "What is 1+1?", RightAnswer: "2", WrongAnswers: []string{"1", "54"}}
}

func printHello() {
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

func readJSON(jsonFile string) []Question {
	file, _ := ioutil.ReadFile(jsonFile)

	var data []Question

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func getUserInput(stdin io.Reader) (string, error) {
	reader := bufio.NewReader(stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, "\n")
	return text, nil
}

func randomizeAnswers(answers []string) []string {
	//TODO: add Seed to truly randomize later
	// See https://yourbasic.org/golang/shuffle-slice-array/
	rand.Seed(0) // predictable shuffling
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	return answers
}

func formatQuestion(testQuestion Question) string {
	var answerOptions = testQuestion.WrongAnswers
	answerOptions = append(answerOptions, testQuestion.RightAnswer)
	randomizedAnswers := randomizeAnswers(answerOptions)
	return fmt.Sprintf(
		"Question: %s\n1: %s\nX: %s\n2: %s\nAnswer: ",
		testQuestion.Question,
		randomizedAnswers[0],
		randomizedAnswers[1],
		randomizedAnswers[2])
}

func main() {

}
