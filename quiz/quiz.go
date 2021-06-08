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
	"net/url"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type TriviaObject struct {
	BaseURL    string `yaml:"base_url"`
	Amount     string `yaml:"amount"`
	Category   string `yaml:"category"`
	Difficulty string `yaml:"difficulty"`
}

type Configuration struct {
	QuestionFile string `yaml:"question_file"`
	//TriviaURL    string `yaml:"trivia_url"`
	Trivia TriviaObject
}

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
	ReadConfigurationFromYAML(yamlFile string) Configuration
	GetQuestions(configuration Configuration) []Question
	GetAnswerMap(question Question, randomizeAnswers bool) map[string]string
	GetUserInput(stdin io.Reader) (string, error)
	FormatQuestion(question Question, answerMap map[string]string) string
	Verify(question Question, answerMap map[string]string, userInput string) (bool, error)
	FormatResult(correctAnswers int, numberQuestions int) string
}

type Quiz struct {
	questions []Question
}

func (quiz *Quiz) GetQuestions(configuration Configuration) []Question {
	triviaUrl, err := quiz.createTriviaURL(configuration)
	if err != nil {
		fmt.Printf("Failed to build URL: %s\n", err.Error())
	}
	// TODO: Avoid always reading from URL? (if createTriviaURL failed, for instance)
	question, err := quiz.readQuestionsFromURL(triviaUrl)
	if err != nil {
		fmt.Printf("Failed to read OpenTrivia, reading from file: %s\n", err.Error())
		question = quiz.readQuestionsFromJSON(configuration.QuestionFile)
	}

	return question
}

func (quiz *Quiz) readQuestionsFromJSON(jsonFile string) []Question {
	file, _ := ioutil.ReadFile(jsonFile)

	var data []Question

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func (quiz *Quiz) ReadConfigurationFromYAML(yamlFile string) Configuration {
	file, _ := ioutil.ReadFile(yamlFile)

	var data Configuration

	_ = yaml.Unmarshal([]byte(file), &data)

	return data
}

func (quiz *Quiz) createTriviaURL(configuration Configuration) (string, error) {
	base := configuration.Trivia.BaseURL
	amount := configuration.Trivia.Amount
	if base == "" || amount == "" {
		err := fmt.Errorf("Mandatory configurations 'base_url' or/and 'amount' missing")
		fmt.Println("Error:", err.Error())
		return "", err
	}

	triviaUrl, err := url.Parse(base)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return "", err

	}
	if triviaUrl.Scheme == "" || triviaUrl.Host == "" {
		err := fmt.Errorf("base_url is missing scheme or host")
		fmt.Println("Error:", err.Error())
		return "", err
	}
	params := url.Values{}
	params.Add("amount", amount)
	params.Add("type", "multiple")

	category := configuration.Trivia.Category
	difficulty := configuration.Trivia.Difficulty
	if category != "" {
		params.Add("category", category)
	}
	if difficulty != "" {
		params.Add("difficulty", difficulty)
	}

	triviaUrl.RawQuery = params.Encode()

	return triviaUrl.String(), nil
}

func (quiz *Quiz) readQuestionsFromURL(url string) ([]Question, error) {
	// GET OpenTrivia questions

	resp, err := http.Get(url)
	var questions []Question

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return questions, err
	}

	if resp.StatusCode != 200 {
		err := fmt.Errorf("Http request not OK: %s", resp.Status)
		fmt.Printf("ERROR: %s\n", err.Error())
		return questions, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var openTriviaResponse OpenTriviaResponse

	err = json.Unmarshal(body, &openTriviaResponse)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return questions, err
	}

	if len(openTriviaResponse.Results) == 0 {
		err := fmt.Errorf("Unable to resolve response into question(s): %s", body)
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
		"\nQuestion: %s\n1: %s\n2: %s\n3: %s\n4: %s\nAnswer: ",
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

	shuffledAnswers := answers
	rand.Seed(seed)
	rand.Shuffle(len(shuffledAnswers), func(i, j int) { shuffledAnswers[i], shuffledAnswers[j] = shuffledAnswers[j], shuffledAnswers[i] })

	return shuffledAnswers
}
