package services_test

import (
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ultranaco/gocurl/services"
)

func TestRunCurl(t *testing.T) {
	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 5 * time.Minute,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          2000,
			MaxIdleConnsPerHost:   400,
			IdleConnTimeout:       5 * time.Minute,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
		},
	}
	logDefault := log.Default()

	uri := "https://postman-echo.com/post"

	rampageLimit := 3
	concurrency := 5

	runner := services.CurlRunner{}
	stats := services.Stats{}

	rampage := time.Duration(1) * time.Minute
	timeout := time.Duration(5) * time.Minute

	headers := make(map[string]string)

	headers["Content-Type"] = "application/json"

	contents, _ := os.ReadFile("../testdata/payload.json")

	runner.New(contents, rampage, rampageLimit, concurrency,
		uri, "POST", httpClient, logDefault, &stats, headers)

	runner.Run()

	runner.HeartBeat(runner.Start.Add(timeout))
}
