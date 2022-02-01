package services

type StatsRampage struct {
	ramapgeID    int
	concurrecy   int
	requestTotal int
	errorTotal   int
	errors       map[string]int
	requestCodes map[int]int
}

func (service *StatsRampage) New(rampageID, concurrency int) {
	service.ramapgeID = rampageID
	service.concurrecy = concurrency
	service.errors = make(map[string]int)
	service.requestCodes = make(map[int]int)
}

func (service StatsRampage) HitRequest(statusCode int) StatsRampage {
	service.requestTotal++
	service.requestCodes[statusCode]++
	return service
}

func (service StatsRampage) HitError(err error) StatsRampage {
	service.errorTotal++
	service.errors[err.Error()]++
	return service
}
