package quiz

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

type Quiz struct {
	questions []Question
}

func (quiz *Quiz) readQuestionsFromJSON(jsonFile string) []Question {
	file, _ := ioutil.ReadFile(jsonFile)

	var data []Question

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func (quiz *Quiz) getUserInput(stdin io.Reader) (string, error) {
	reader := bufio.NewReader(stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	text = strings.TrimSuffix(text, "\n")
	return text, nil
}

func (quiz *Quiz) randomizeAnswers(answers []string) []string {
	//TODO: add Seed to truly randomize later
	// See https://yourbasic.org/golang/shuffle-slice-array/
	rand.Seed(0) // predictable shuffling
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	return answers
}

func (quiz *Quiz) formatQuestion(testQuestion Question) (string, []string) {
	var answerOptions = testQuestion.WrongAnswers
	answerOptions = append(answerOptions, testQuestion.RightAnswer)
	randomizedAnswers := quiz.randomizeAnswers(answerOptions)
	questionAndAnswers := fmt.Sprintf(
		"Question: %s\n1: %s\nX: %s\n2: %s\nAnswer: ",
		testQuestion.Question,
		randomizedAnswers[0],
		randomizedAnswers[1],
		randomizedAnswers[2])

	return questionAndAnswers, randomizedAnswers
}

func (quiz *Quiz) getAnswerMap(question Question) map[string]string {
	var answerOptions = question.WrongAnswers
	answerOptions = append(answerOptions, question.RightAnswer)
	randomizedAnswers := quiz.randomizeAnswers(answerOptions)

	return map[string]string{
		"1": randomizedAnswers[0],
		"x": randomizedAnswers[1],
		"2": randomizedAnswers[2],
	}

}

// This function verifies that the answer is correct
func (quiz *Quiz) verify(testQuestion Question, answer string) (bool, error) {

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
