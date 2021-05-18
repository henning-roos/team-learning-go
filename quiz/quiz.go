package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Question struct {
	Question     string    `json:"question"`
	RightAnswer  string    `json:"correct_answer"`
	WrongAnswers [3]string `json:"incorrect_answers"`
}

type OpenTriviaResponse struct {
	ResponseCode int        `json:"response_code"`
	Results      []Question `json:"results"`
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

func (quiz *Quiz) ReadQuestionsFromURL(url string) ([]Question, error) {
	// GET OpenTrivia questions

	resp, err := http.Get(url)
	var questions []Question

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return questions, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Invalid response code %s\n", resp.Status)
		return questions, fmt.Errorf("Http request not OK: %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var openTriviaResponse OpenTriviaResponse

	err = json.Unmarshal(body, &openTriviaResponse)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return questions, err
	}

	for i, question := range openTriviaResponse.Results {
		question.Question = html.UnescapeString(question.Question)
		question.RightAnswer = html.UnescapeString(question.RightAnswer)
		question.WrongAnswers[0] = html.UnescapeString(question.WrongAnswers[0])
		question.WrongAnswers[1] = html.UnescapeString(question.WrongAnswers[1])
		question.WrongAnswers[2] = html.UnescapeString(question.WrongAnswers[2])
		openTriviaResponse.Results[i] = question
	}

	return openTriviaResponse.Results, nil
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
		"Question: %s\n1: %s\n2: %s\n3: %s\n4: %s\nAnswer: ",
		question.Question,
		answerMap["1"],
		answerMap["2"],
		answerMap["3"],
		answerMap["4"])

	return questionAndAnswers
}

func (quiz *Quiz) GetAnswerMap(question Question, randomizeSeed bool) map[string]string {
	var answerOptions []string
	answerOptions = append(answerOptions, question.WrongAnswers[0])
	answerOptions = append(answerOptions, question.WrongAnswers[1])
	answerOptions = append(answerOptions, question.WrongAnswers[2])
	answerOptions = append(answerOptions, question.RightAnswer)
	randomizedAnswers := quiz.randomizeAnswers(answerOptions, randomizeSeed)

	return map[string]string{
		"1": randomizedAnswers[0],
		"2": randomizedAnswers[1],
		"3": randomizedAnswers[2],
		"4": randomizedAnswers[3],
	}
}

// This function verifies that the answer is correct
func (quiz *Quiz) Verify(question Question, answerMap map[string]string, userInput string) (bool, error) {

	//Assume userInput is 1, 2, 3 or 4
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

// Function that format the result printout to console
func (quiz *Quiz) FormatResult(numberCorrectAnswers int, numberQuestions int) string {
	resultString := fmt.Sprintf("You got %d of %d correct answers.", numberCorrectAnswers, numberQuestions)

	return resultString
}

func (quiz *Quiz) randomizeAnswers(answers []string, randomizeSeed bool) []string {

	var seed int64 = 0 // Default to 0 for deterministic testing

	if randomizeSeed {
		// See https://yourbasic.org/golang/shuffle-slice-array/
		seed = time.Now().UnixNano()
	}

	shuffeledAnswers := answers
	rand.Seed(seed)
	rand.Shuffle(len(shuffeledAnswers), func(i, j int) { shuffeledAnswers[i], shuffeledAnswers[j] = shuffeledAnswers[j], shuffeledAnswers[i] })

	return shuffeledAnswers
}
