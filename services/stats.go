package services

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	mutex                       *sync.RWMutex
	format                      Format
	errorTotal                  int
	errors                      map[string]int
	errorTimeout                int
	errorTimeoutDial            int
	errorTimeoutTLS             int
	errorTimeoutResopnseHeaders int
	errorTimeoutIdleConn        int
	requestTotal                int
	requestCodes                map[int]int
	concurrency                 int
	rampageIteration            int
}

func (service *Stats) New() {
	service.mutex = &sync.RWMutex{}
	service.requestCodes = make(map[int]int)
	service.errors = make(map[string]int)
}

func (service *Stats) HitRequest(statusCode int) {
	go func(scode int) {
		service.mutex.Lock()
		service.requestTotal++
		service.requestCodes[statusCode]++
		service.mutex.Unlock()
	}(statusCode)
}

func (service *Stats) HitError(err error) {
	go func(e error) {
		service.mutex.Lock()
		service.errorTotal++
		service.errors[err.Error()]++
		service.mutex.Unlock()
	}(err)
}

func (service Stats) Print(start, end, current time.Time) {
	progress := service.format.ProgressFromDateTime(start, end, current)
	progressBar := service.format.TimeProgressLoad(progress)
	service.mutex.RLock()
	fmt.Print("\033[u\033[K")
	fmt.Printf(`
Time:		%v
Progress:	`+"\033[32m"+`%s`+"\033[0m"+`
Request:	%v
Error:		%v`, current, progressBar, service.requestTotal, service.errorTotal)
	if progress >= 100 {
		fmt.Print("\n\n")
		for code, value := range service.requestCodes {
			fmt.Printf("Status Codes:	%v - [%v]\n", code, value)
		}
		fmt.Print("\n")

		for code, value := range service.errors {
			fmt.Printf("Errors:		[%v] \033[37;41m%s\033[0m\n", value, code)
		}
		fmt.Print("\n")

	}
	service.mutex.RUnlock()
}
