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

type QuizInterface interface {
	ReadQuestionsFromJSON(jsonFile string) []Question
	GetAnswerMap(question Question) map[string]string
	GetUserInput(stdin io.Reader) (string, error)
	FormatQuestion(question Question, answerMap map[string]string) string
	Verify(question Question, answerMap map[string]string, userInput string) (bool, error)
}

type Quiz struct {
	questions []Question
}

func (quiz *Quiz) ReadQuestionsFromJSON(jsonFile string) []Question {
	file, _ := ioutil.ReadFile(jsonFile)

	var data []Question

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func (quiz *Quiz) GetUserInput(stdin io.Reader) (string, error) {
	reader := bufio.NewReader(stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")
	return text, nil
}

func (quiz *Quiz) FormatQuestion(question Question, answerMap map[string]string) string {
	questionAndAnswers := fmt.Sprintf(
		"Question: %s\n1: %s\nX: %s\n2: %s\nAnswer: ",
		question.Question,
		answerMap["1"],
		answerMap["X"],
		answerMap["2"])

	return questionAndAnswers
}

func (quiz *Quiz) GetAnswerMap(question Question) map[string]string {
	var answerOptions = question.WrongAnswers
	answerOptions = append(answerOptions, question.RightAnswer)
	randomizedAnswers := quiz.randomizeAnswers(answerOptions)

	return map[string]string{
		"1": randomizedAnswers[0],
		"X": randomizedAnswers[1],
		"2": randomizedAnswers[2],
	}

}

// This function verifies that the answer is correct
func (quiz *Quiz) Verify(question Question, answerMap map[string]string, userInput string) (bool, error) {

	//Assume userInput is 1, X or 2
	userAnswer := answerMap[userInput]

	if userAnswer == question.RightAnswer {
		return true, nil
	}

	for _, value := range question.WrongAnswers {
		if userAnswer == value {
			return false, nil
		}
	}

	return false, fmt.Errorf("The specified answer is invalid answer: %s", userInput)
}

func (quiz *Quiz) randomizeAnswers(answers []string) []string {
	//TODO: add Seed to truly randomize later
	// See https://yourbasic.org/golang/shuffle-slice-array/
	rand.Seed(0) // predictable shuffling
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	return answers
}
