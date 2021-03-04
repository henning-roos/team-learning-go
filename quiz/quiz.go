package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Question struct {
	Question     string   `json:"question"`
	RightAnswer  string   `json:"rightAnswer"`
	WrongAnswers []string `json:"wrongAnswers"`
}

type QuizInterface interface {
	ReadQuestionsFromJSON(jsonFile string) []Question
	GetAnswerMap(question Question, randomizeAnswers bool) map[string]string
	GetUserInput(stdin io.Reader) (string, error)
	FormatQuestion(question Question, answerMap map[string]string) string
	Verify(question Question, answerMap map[string]string, userInput string) (bool, error)
	FormatResult(correctAnswers int, numberQuestions int) string
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

func (quiz *Quiz) GetAnswerMap(question Question, randomizeSeed bool) map[string]string {
	var answerOptions = question.WrongAnswers
	answerOptions = append(answerOptions, question.RightAnswer)
	randomizedAnswers := quiz.randomizeAnswers(answerOptions, randomizeSeed)

	return map[string]string{
		"1": randomizedAnswers[0],
		"X": randomizedAnswers[1],
		"2": randomizedAnswers[2],
	}

}

// This function verifies that the answer is correct
func (quiz *Quiz) Verify(question Question, answerMap map[string]string, userInput string) (bool, error) {

	// fmt.Printf("QUESTION: %v\n", question)
	// fmt.Printf("ANSWERMAP: %v\n", answerMap)

	//Assume userInput is 1, x, X or 2
	// fmt.Println("Userinput before toUpper:|" + userInput + "|")
	userAnswer := answerMap[strings.ToUpper(userInput)]
	// fmt.Println("converted userinput to Uppercase" + strings.ToUpper(userInput))
	// fmt.Println("usrAnswer from map is: " + userAnswer)
	// fmt.Println(answerMap)
	// fmt.Println("right answer from map: " + question.RightAnswer)
	// fmt.Printf("wrong answers from map: %v \n", question.WrongAnswers)

	if userAnswer == question.RightAnswer {
		return true, nil
	}

	for _, value := range question.WrongAnswers {
		fmt.Printf("Compare Useranswer %s with wront answer %s \n", userAnswer, value)
		if userAnswer == value {
			return false, nil
		}
	}

	return false, fmt.Errorf("The specified answer is invalid answer: %s", userInput)
}

// Function that format the result printout to console
func (quiz *Quiz) FormatResult(numberCorrectAnswers int, numberQuestions int) string {
	resultString := fmt.Sprintf("You got %d of %d correct answers.", numberCorrectAnswers, numberQuestions)

	return resultString
}

func (quiz *Quiz) randomizeAnswers(answers []string, randomizeSeed bool) []string {

	var seed int64 = 12 // Default to 0 for deterministic testing

	if randomizeSeed {
		// See https://yourbasic.org/golang/shuffle-slice-array/
		seed = time.Now().UnixNano()
	}

	rand.Seed(seed)
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	return answers
}
