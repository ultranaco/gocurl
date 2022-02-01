package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/ultranaco/gocurl/services"
)

func TestParsePayload(t *testing.T) {
	parser := services.ParseArgs{}
	contents, _ := parser.ParsePayload("@/home/alex/dev/golang/gocurl/testdata/payload.json")

	t.Log(string(contents))
}

func TestGetHeadersHappyPath(t *testing.T) {
	parser := services.ParseArgs{}

	headres, err := parser.GetHeaders("Content-Type : application/json", "User-Agent:Mozilla")

	if err != nil {
		t.Logf("error while parsing header %v", err)
		t.Fail()
		return
	}

	t.Log(headres)
}

func TestGetLogger(t *testing.T) {
	parser := services.ParseArgs{}

	log.Println("starting test")

	logger, err := parser.GetHttpResponseLogger("/dev/null", 0)

	if err != nil {
		t.Logf("error while parsing to retrieve logger %v", err)
		t.Fail()
		return
	}

	logger.Println("nothig happens")

	logger, err = parser.GetHttpResponseLogger("", 10)

	if err != nil {
		t.Logf("error while parsing to retrieve logger %v", err)
		t.Fail()
		return
	}

	logger.Println("nothig happens")

	logger, err = parser.GetHttpResponseLogger("", 0)

	if err != nil {
		t.Logf("error while parsing to retrieve logger %v", err)
		t.Fail()
		return
	}

	logger.Println("happens somethig output in console")
	log.Println("happens somethig output in console")

	logger, err = parser.GetHttpResponseLogger("../testdata/logs.txt", 10)

	if err != nil {
		t.Logf("error while parsing to retrieve logger %v", err)
		t.Fail()
		return
	}

	logger.Println("happens something in file log")
	log.Println("again somethig output in console")

}

func TestParseTimeDefault(t *testing.T) {

	parser := services.ParseArgs{}
	duration, _ := parser.GetTimeDuration("11")

	if duration != (time.Millisecond * 11) {
		t.Logf("%v duration 11ms fail", duration)
		t.Fail()
		return
	}

	t.Logf("%v duration 12ms succes", duration)
}

func TestParseTimeWithUnits(t *testing.T) {

	parser := services.ParseArgs{}
	duration, _ := parser.GetTimeDuration("12ms")

	if duration != (time.Millisecond * 12) {
		t.Logf("%v duration 12 miliseconds fail\n", duration)
		t.Fail()
		return
	}

	t.Logf("%v duration 12 miliseconds succes\n", duration)

	duration, _ = parser.GetTimeDuration("13s")

	if duration != (time.Second * 13) {
		t.Logf("%v duration 13 seconds fail\n", duration)
		t.Fail()
		return
	}

	t.Logf("%v duration 13 seconds succes\n", duration)

	duration, _ = parser.GetTimeDuration("14m")

	if duration != (time.Minute * 14) {
		t.Logf("%v duration 14 minutes fail\n", duration)
		t.Fail()
		return
	}

	t.Logf("%v duration 14 minutes succes\n", duration)

	duration, _ = parser.GetTimeDuration("15h")

	if duration != (time.Hour * 15) {
		t.Logf("%v duration 15 hours fail\n", duration)
		t.Fail()
		return
	}

	t.Logf("%v duration 15 hours succes\n", duration)
}

func TestHandleErrors(t *testing.T) {
	parser := services.ParseArgs{}

	_, err := parser.GetTimeDuration("")

	if err == nil {
		t.Logf("%v\n", err)
		t.Fail()
		return
	}

	t.Logf("empty time error handle successfully\n")

	_, err = parser.GetTimeDuration("3y")

	if err == nil {
		t.Logf("%v\n", err)
		t.Fail()
		return
	}

	t.Logf("unit not allowed error handle successfully\n")

	_, err = parser.GetTimeDuration("f1m")

	if err == nil {
		t.Logf("%v\n", err)
		t.Fail()
		return
	}

	t.Logf("pattern not match error handle successfully\n")
}
