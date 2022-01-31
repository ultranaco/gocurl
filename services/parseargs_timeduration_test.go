package services_test

import (
	"testing"
	"time"

	"github.com/ultranaco/gocurl/services"
)

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
