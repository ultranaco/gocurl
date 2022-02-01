package services

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	mutex            *sync.RWMutex
	format           Format
	errorTotal       int
	requestTotal     int
	concurrency      int
	rampageIteration int
	requestCodes     map[int]int
	errors           map[string]int
	rampageStats     map[int]StatsRampage
}

func (service *Stats) New() {
	service.mutex = &sync.RWMutex{}
	service.requestCodes = make(map[int]int)
	service.errors = make(map[string]int)
	service.rampageStats = make(map[int]StatsRampage)
}

func (service *Stats) HitRequest(statusCode int) {
	go func(scode int) {
		service.mutex.Lock()

		service.requestTotal++
		service.requestCodes[statusCode]++
		rampage := service.rampageStats[service.rampageIteration]
		rampage = rampage.HitRequest(statusCode)
		service.rampageStats[service.rampageIteration] = rampage

		service.mutex.Unlock()
	}(statusCode)
}

func (service *Stats) HitError(err error) {
	go func(e error) {

		service.mutex.Lock()

		service.errorTotal++
		service.errors[err.Error()]++
		rampage := service.rampageStats[service.rampageIteration]
		rampage = rampage.HitError(err)
		service.rampageStats[service.rampageIteration] = rampage

		service.mutex.Unlock()

	}(err)
}

func (service *Stats) HitRampage(ramapage, concurrency int) {
	go func(ramp, conc int) {

		service.mutex.Lock()
		service.rampageIteration = ramapage
		service.concurrency += concurrency

		rampage := StatsRampage{}
		rampage.New(ramapage, service.concurrency)
		service.rampageStats[ramapage] = rampage

		service.mutex.Unlock()

	}(ramapage, concurrency)
}

func (service Stats) Print(start, end, current time.Time) {
	progress := service.format.ProgressFromDateTime(start, end, current)
	progressBar := service.format.TimeProgressLoad(progress)
	service.mutex.RLock()
	fmt.Print("\033[u\033[K")
	fmt.Printf(`
time:		%v
progress:	`+"\033[32m"+`%s`+"\033[0m"+`
rampage:	`+"\033[96m"+`[%v]`+"\033[0m"+`
request:	`+"\033[96m"+`[%v]`+"\033[0m"+`
error:		`+"\033[96m"+`[%v]`+"\033[0m"+`
concurrency:	`+"\033[96m"+`[%v]`+"\033[0m", current, progressBar,
		service.rampageIteration, service.requestTotal, service.errorTotal,
		service.concurrency)

	service.mutex.RUnlock()
}

func (service Stats) PrintSummary() {
	fmt.Print("\n\n")
	fmt.Print(
		"\033[93m____________________Summary Stats____________________\033[0m\n\n")
	for code, value := range service.requestCodes {
		fmt.Printf("status_code:	[%v]	%v\n\n", value, code)
	}

	for code, value := range service.errors {
		fmt.Printf("error_code:	[%v]	\033[37;41m%s\033[0m\n\n", value, code)
	}

	fmt.Print(
		"\033[93m____________________Rampage Stats____________________\033[0m\n\n")
	for rampID := 0; rampID <= service.rampageIteration; rampID++ {

		rampage := service.rampageStats[rampID]

		fmt.Printf(`rampageID:	[%v]
concurrency:	[%v]
request:	[%v]
error:		[%v]

`, rampID, rampage.concurrecy, rampage.requestTotal, rampage.errorTotal)
	}
}
