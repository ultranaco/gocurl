package services

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CurlRunner struct {
	payload        []byte
	rampage        time.Duration
	rampageLimit   int
	url            string
	concurrent     int
	httpClient     *http.Client
	responseLogger *log.Logger
	stats          *Stats
	requestMethod  string
	headers        map[string]string
	Start          time.Time
}

func (runner *CurlRunner) New(
	payload []byte, rampage time.Duration, rampageLimit, concurrent int,
	url, requestMethod string, httpClient *http.Client,
	responseLogger *log.Logger, stats *Stats, headers map[string]string) {
	runner.payload = payload
	runner.rampage = rampage
	runner.url = url
	runner.concurrent = concurrent
	runner.rampageLimit = rampageLimit
	runner.httpClient = httpClient
	runner.stats = stats
	runner.requestMethod = requestMethod
	runner.headers = headers
	runner.responseLogger = responseLogger
	runner.Start = time.Now()
}

func (runner CurlRunner) Run() {
	if runner.concurrent <= 0 {
		runner.request(false)
		return
	}

	runner.rampageRun(func() {
		for iteraton := 0; iteraton < runner.concurrent; iteraton++ {
			go runner.request(true)
		}
	})
}

func (runner CurlRunner) HeartBeat(timeout time.Time) {
	if runner.concurrent <= 0 {
		return
	}
	ticker := time.NewTicker(1 * time.Second)
	for tick := range ticker.C {
		runner.stats.Print(runner.Start, timeout, tick)
		if tick.Equal(timeout) || tick.After(timeout) {
			return
		}
	}
}

func (runner CurlRunner) rampageRun(action func()) {
	if runner.rampage == time.Duration(0) {
		action()
		return
	}

	go func() {
		iteration := 0
		for {
			runner.stats.HitRampage(iteration, runner.concurrent)
			action()
			if iteration == runner.rampageLimit {
				break
			}

			iteration++
			time.Sleep(runner.rampage)
		}
	}()
}

func (runner CurlRunner) request(infinity bool) {
	for iterator := 0; iterator < 1 || infinity; iterator++ {
		httpRequest, err := http.NewRequest(
			http.MethodPost, runner.url, bytes.NewReader(runner.payload))

		if value, ok := runner.headers["Content-Type"]; ok {
			httpRequest.Header.Add("Content-Type", value)
		}

		if err != nil {
			continue
		}

		contents, err := runner.doRequest(httpRequest)

		if err != nil {
			continue
		}

		runner.responseLogger.Println(string(contents))
	}
}

func (service CurlRunner) doRequest(request *http.Request) ([]byte, error) {
	response, err := service.httpClient.Do(request)
	statusCode := 0

	if response != nil {
		statusCode = response.StatusCode
	}

	service.stats.HitRequest(statusCode)

	if err != nil {
		service.stats.HitError(err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		service.stats.HitError(err)
		return nil, err
	}

	return body, nil
}
