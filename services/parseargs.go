package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var timePattern = regexp.MustCompile(`^(\d+)(ms|[smh])?$`)
var patternHeader = regexp.MustCompile(`^([a-zA-Z0-9\-_]+?)\s*?\:\s*?([a-zA-Z].+)$`)

type ParseArgs struct {
}

func (parser ParseArgs) GetTimeDuration(rawElapsed string) (time.Duration, error) {
	matches := timePattern.FindAllSubmatch([]byte(rawElapsed), 1)
	var duration time.Duration
	var err error

	if len(matches) <= 0 {
		return time.Duration(0), errors.New("ERROR: input a valid format [digit+(ms|s|m|h)?]")
	}

	matchTime := matches[0]

	rawConstant := matchTime[1]
	rawUnit := matchTime[2]

	constant, _ := strconv.Atoi(string(rawConstant))

	switch string(rawUnit) {
	case "h":
		duration = time.Duration(constant) * time.Hour
	case "m":
		duration = time.Duration(constant) * time.Minute
	case "s":
		duration = time.Duration(constant) * time.Second
	case "", "ms":
		duration = time.Duration(constant) * time.Millisecond
	default:
		err = errors.New("ERROR: unit time not recognized you should use (ms|s|m|h)")
	}

	return duration, err
}

func (parser ParseArgs) ParsePayload(input string) ([]byte, error) {

	if input == "" {
		return nil, nil
	}

	var isFile bool

	if input[0:1] == "@" {
		isFile = true
	}

	if isFile {
		filePath := input[1:]
		return os.ReadFile(filePath)
	}

	return []byte(input), nil
}

func (parser ParseArgs) GetMehod(inputMethod string) (string, error) {
	if inputMethod == "" {
		return "", errors.New("ERROR: request method can not be empty")
	}

	input := strings.ToUpper(inputMethod)

	switch input {
	case "GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH":
		return input, nil
	default:
		return input, fmt.Errorf("ERROR: unrecognized request method [%s]", inputMethod)
	}
}

func (parser ParseArgs) GetHeaders(inputHeaders ...string) (map[string]string, error) {
	headers := make(map[string]string)

	for _, inputHeader := range inputHeaders {
		if patternHeader.Match([]byte(inputHeader)) {
			matches := patternHeader.FindAllSubmatch([]byte(inputHeader), 1)
			headers[string(matches[0][1])] = string(matches[0][2])
		} else {
			return headers, fmt.Errorf("ERROR: format header malformed \"Header: Value\", [%s]", inputHeader)
		}
	}

	return headers, nil
}

func (parser ParseArgs) GetHttpResponseLogger(logOutput string, concurrency int) (*log.Logger, error) {
	if (concurrency > 0 && logOutput == "") || logOutput == "/dev/null" {
		logDefault := log.New(os.Stderr, "", log.LstdFlags)
		logDefault.SetOutput(ioutil.Discard)
		return logDefault, nil
	}

	if logOutput == "" || logOutput == "std" {
		logDefault := log.New(os.Stderr, "", log.LstdFlags)
		return logDefault, nil
	}

	file, err := os.OpenFile(logOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return nil, err
	}

	logger := log.New(file, "", log.Ldate|log.Ltime)

	return logger, nil
}
