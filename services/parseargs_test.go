package services_test

import (
	"log"
	"testing"

	"github.com/ultranaco/gocurl/services"
)

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
