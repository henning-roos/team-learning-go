package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testConfiguration = Configuration{
	QuestionFile: "Path.To.Some.File",
	TriviaURL:    "url://to.trivia",
}

var testQuestion = Question{
	Question:     "Which language is this written in?",
	RightAnswer:  "Go",
	WrongAnswers: [3]string{"Python", "Java", "Ruby"},
}

var testQuestion2 = Question{
	Question:     "What is blue and yellow together? (using watercolors)",
	RightAnswer:  "Green",
	WrongAnswers: [3]string{"Red", "Black", "Pink"},
}

var testAnswerMap = map[string]string{
	"1": "Ruby",
	"2": "Java",
	"3": "Python",
	"4": "Go",
}

var testAnswerMap2 = map[string]string{
	"1": "Pink",
	"2": "Black",
	"3": "Red",
	"4": "Green",
}

var quiz = Quiz{}

func TestGetQuestions(t *testing.T) {
	okResponse := func(res http.ResponseWriter, req *http.Request) {
		jsonData := `{"response_code":0,"results":[{"category":"Entertainment: Video Games","type":"multiple","difficulty":"medium","question":"In &quot;Call Of Duty: Zombies&quot;, which map features the &quot;Fly Trap&quot; easter egg?","correct_answer":"Der Riese","incorrect_answers":["Tranzit","Call Of The Dead","Shi No Numa"]}]}`
		res.Write([]byte(jsonData))
	}

	testServer := httptest.NewServer(http.HandlerFunc(okResponse))
	defer func() { testServer.Close() }()

	var testConfiguration = Configuration{
		QuestionFile: "questions.json",
		TriviaURL:    testServer.URL,
	}

	questions := quiz.GetQuestions(testConfiguration)
	assert.Equal(t, 1, len(questions))
}

func TestGetQuestionsBadUrl(t *testing.T) {
	badResponse := func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(503)
	}

	testServer := httptest.NewServer(http.HandlerFunc(badResponse))
	defer func() { testServer.Close() }()

	var testConfiguration = Configuration{
		QuestionFile: "questions.json",
		TriviaURL:    testServer.URL,
	}

	questions := quiz.GetQuestions(testConfiguration)
	assert.Equal(t, 3, len(questions))
}

func TestReadQuestionsFromJSON(t *testing.T) {
	questions := quiz.readQuestionsFromJSON("questions.json")
	assert.Equal(t, "What is blue and yellow together? (using watercolors)", questions[0].Question)
	assert.Equal(t, "Green", questions[0].RightAnswer)
	assert.Equal(t, "Red", questions[0].WrongAnswers[0])
	assert.Equal(t, 3, len(questions))
}

func TestReadConfigurationFromYAML(t *testing.T) {
	configuration := quiz.ReadConfigurationFromYAML("config.yaml")
	assert.NotEmpty(t, configuration.TriviaURL)
	assert.NotEmpty(t, configuration.QuestionFile)
	assert.NotEmpty(t, configuration.Trivia.BaseURL)
	assert.NotEmpty(t, configuration.Trivia.Amount)
}

func TestCreateTriviaURL(t *testing.T) {
	var trivia = TriviaObject{
		BaseURL:    "trivia.com/api",
		Amount:     "10",
		Category:   "s",
		Difficulty: "s",
	}
	var testConfiguration = Configuration{
		QuestionFile: "questions.json",
		TriviaURL:    "trivia.com/api",
		Trivia:       trivia,
	}

	triviaURL, _ := quiz.createTriviaURL(testConfiguration)
	assert.Equal(t, "trivia.com/api?amount=10&category=s&difficulty=s&type=multiple", triviaURL)
}

func TestCreateTriviaURLMissingMandatory(t *testing.T) {
	var trivia = TriviaObject{
		BaseURL: "trivia.com/api",
	}
	var testConfiguration = Configuration{
		QuestionFile: "questions.json",
		TriviaURL:    "trivia.com/api",
		Trivia:       trivia,
	}

	_, err := quiz.createTriviaURL(testConfiguration)
	assert.Error(t, err)

}

func TestCreateTriviaURLMissingHost(t *testing.T) {
	var trivia = TriviaObject{
		BaseURL: "this is not an url",
		Amount:  "10",
	}
	var testConfiguration = Configuration{
		QuestionFile: "questions.json",
		TriviaURL:    "trivia.com/api",
		Trivia:       trivia,
	}

	_, err := quiz.createTriviaURL(testConfiguration)
	assert.Equal(t, err.Error(), "base_url is missing scheme or host")
}

func TestCreateTriviaURLMissingSchema(t *testing.T) {
	var trivia = TriviaObject{
		BaseURL: "this is not an url",
		Amount:  "10",
	}
	var testConfiguration = Configuration{
		QuestionFile: "questions.json",
		TriviaURL:    "trivia.com/api",
		Trivia:       trivia,
	}

	_, err := quiz.createTriviaURL(testConfiguration)
	assert.Equal(t, err.Error(), "base_url is missing scheme or host")
}

func TestReadQuestionsFromURL(t *testing.T) {
	okResponse := func(res http.ResponseWriter, req *http.Request) {
		jsonData := `{"response_code":0,"results":[{"category":"Entertainment: Video Games","type":"multiple","difficulty":"medium","question":"In &quot;Call Of Duty: Zombies&quot;, which map features the &quot;Fly Trap&quot; easter egg?","correct_answer":"Der Riese","incorrect_answers":["Tranzit","Call Of The Dead","Shi No Numa"]}]}`
		res.Write([]byte(jsonData))
	}

	testServer := httptest.NewServer(http.HandlerFunc(okResponse))
	defer func() { testServer.Close() }()

	questions, _ := quiz.readQuestionsFromURL(testServer.URL)
	questionText := `In "Call Of Duty: Zombies", which map features the "Fly Trap" easter egg?`
	assert.Equal(t, questionText, questions[0].Question)
	assert.Equal(t, "Der Riese", questions[0].RightAnswer)
	assert.Equal(t, "Tranzit", questions[0].WrongAnswers[0])
	assert.Equal(t, 1, len(questions))
}

func TestReadQuestionsFromURLWithoutProtocol(t *testing.T) {
	url := "google.se"

	_, err := quiz.readQuestionsFromURL(url)
	assert.Error(t, err)
}

func TestReadQuestionsFromURLWithBadResponse(t *testing.T) {
	badResponse := func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(503)
	}

	testServer := httptest.NewServer(http.HandlerFunc(badResponse))
	defer func() { testServer.Close() }()

	_, err := quiz.readQuestionsFromURL(testServer.URL)
	assert.Error(t, err)
}

func TestReadQuestionsFromURLWithNoJSONResponse(t *testing.T) {
	okResponse := func(res http.ResponseWriter, req *http.Request) {
		jsonData := `no json here`
		res.Write([]byte(jsonData))
	}

	testServer := httptest.NewServer(http.HandlerFunc(okResponse))
	defer func() { testServer.Close() }()

	_, err := quiz.readQuestionsFromURL(testServer.URL)
	assert.Error(t, err)
}

func TestReadQuestionsFromURLWithBadJSONResponse(t *testing.T) {
	okResponse := func(res http.ResponseWriter, req *http.Request) {
		jsonData := `{"key": "no question here"}`
		res.Write([]byte(jsonData))
	}

	testServer := httptest.NewServer(http.HandlerFunc(okResponse))
	defer func() { testServer.Close() }()

	_, err := quiz.readQuestionsFromURL(testServer.URL)
	assert.Error(t, err)
}

func TestGetUserInput(t *testing.T) {
	var stdin bytes.Buffer
	var expected string = "2"
	stdin.Write([]byte(expected + "\n"))
	var input, _ = quiz.GetUserInput(&stdin)
	assert.Equal(t, expected, input)
}

func TestGetUserInputWindows(t *testing.T) {
	var stdin bytes.Buffer
	var expected string = "2"
	stdin.Write([]byte(expected + "\r\n"))
	var input, _ = quiz.GetUserInput(&stdin)
	assert.Equal(t, expected, input)
}

func TestGetUserInputError(t *testing.T) {
	var stdin bytes.Buffer
	stdin.Write([]byte("2"))
	var input, err = quiz.GetUserInput(&stdin)
	assert.Error(t, err)
	assert.Equal(t, "", input)
}

func TestRandomizeAnswers(t *testing.T) {
	var answers = []string{"A", "B", "C"}
	var expected = []string{"B", "A", "C"}
	var randomizeAnswers bool = false
	actual := quiz.randomizeAnswers(answers, randomizeAnswers)
	assert.Equal(t, expected, actual)
}

func TestFormatQuestion(t *testing.T) {
	actualQandA := quiz.FormatQuestion(testQuestion, testAnswerMap)
	expectedQandA := "\nQuestion: Which language is this written in?\n" +
		"1: Ruby\n" +
		"2: Java\n" +
		"3: Python\n" +
		"4: Go\n" +
		"Answer: "
	assert.Equal(t, expectedQandA, actualQandA)
}

func TestGetAnswerMap(t *testing.T) {
	var randomizeAnswers bool = false
	question := testQuestion
	actual := quiz.GetAnswerMap(question, randomizeAnswers)
	assert.Equal(t, testAnswerMap, actual)
}

func TestWrongAnswer(t *testing.T) {
	userInput := "1"
	actual, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.False(t, actual)
	assert.Equal(t, nil, err)
}

func TestCorrectAnswer(t *testing.T) {
	userInput := "4"
	actual, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.True(t, actual)
	assert.Equal(t, nil, err)
}

func TestInvalidAnswer(t *testing.T) {
	var userInput string = ""
	_, err := quiz.Verify(testQuestion, testAnswerMap, userInput)
	assert.Error(t, err)
}

func TestFormatResult(t *testing.T) {
	numberCorrectAnswers := 1
	numberQuestions := 2
	formattedResult := quiz.FormatResult(numberCorrectAnswers, numberQuestions)
	expected := "You got 1 of 2 correct answers."

	assert.Equal(t, expected, formattedResult)
}
