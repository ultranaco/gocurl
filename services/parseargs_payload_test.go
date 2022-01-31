package services_test

import (
	"testing"

	"github.com/ultranaco/gocurl/services"
)

func TestParsePayload(t *testing.T) {
	parser := services.ParseArgs{}
	contents, _ := parser.ParsePayload("@/home/alex/dev/golang/gocurl/testdata/payload.json")

	t.Log(string(contents))
}
