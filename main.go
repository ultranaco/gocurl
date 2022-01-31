package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/ultranaco/gocurl/services"
)

func main() {

	httpClient := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 5 * time.Minute,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          65536,
			MaxIdleConnsPerHost:   2000,
			IdleConnTimeout:       5 * time.Minute,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
		},
	}

	var uri string
	var headers map[string]string

	concurrent := flag.Int("c", 0, "concurrency")
	rampageTime := flag.String("rt", "0", "rampage time window")
	rawTimeout := flag.String("t", "1m", "timeout to finish stress test")
	rampageLimit := flag.Int("rl", 3, "rampage limit is the number of times wich concurrency is multiply by rampage iteration")
	rawOutput := flag.String("o", "", "location to save data, if concurrency is active default id /dev/null")

	rawHeader := flag.String("H", "", "add headers \"Header: Value\"")
	rawRequestMethod := flag.String("X", "GET", "request method")
	data := flag.String("d", "", "data payload")

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		uri = args[len(args)-1]
	}

	parser := services.ParseArgs{}
	runner := services.CurlRunner{}
	stats := &services.Stats{}

	stats.New()

	rampage, err := parser.GetTimeDuration(*rampageTime)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	timeout, err := parser.GetTimeDuration(*rawTimeout)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	contents, err := parser.ParsePayload(*data)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	method, err := parser.GetMehod(*rawRequestMethod)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	headers, err = parser.GetHeaders(*rawHeader)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	httpResponseLogger, err := parser.GetHttpResponseLogger(*rawOutput, *concurrent)

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	runner.New(contents, rampage, *rampageLimit, *concurrent,
		uri, method, httpClient, httpResponseLogger, stats, headers)

	log.Printf("starting!\n\n")
	runner.Run()
	fmt.Print("\n\n\n\n\n\n")
	fmt.Print("\033[6A\033[s")
	runner.HeartBeat(runner.Start.Add(timeout))

	log.Printf("completed!")
}
