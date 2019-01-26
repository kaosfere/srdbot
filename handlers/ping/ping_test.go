package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	var body slack.Msg
	var request events.APIGatewayProxyRequest
	response, _ := handleRequest(request)

	assert := assert.New(t)
	assert.Equal(200, response.StatusCode)

	err := json.Unmarshal([]byte(response.Body), &body)
	if !assert.Nil(err) {
		t.Log("Unmarshal failed.  Unable to continue.")
		t.FailNow()
	}

	assert.Equal("ephemeral", body.ResponseType)
	assert.Equal("SRDBot works!", body.Text)
}

func TestErrorResponse(t *testing.T) {
	err := fmt.Errorf("test")
	response := errorResponse(err)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "{\"error\": \"test\"}", response.Body)
}
