package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

func TestPing(t *testing.T) {
	var body slack.Msg
	var request events.APIGatewayProxyRequest
	response, _ := handleRequest(request)

	if response.StatusCode != 200 {
		t.Errorf("Incorrect status code: %d", response.StatusCode)
	}

	err := json.Unmarshal([]byte(response.Body), &body)
	if err != nil {
		t.Errorf("Error unmarshalling response: %s", err)

	}
	if body.ResponseType != "ephemeral" {
		t.Errorf("Incorrect ResponseType: %s", body.ResponseType)
	}

	if body.Text != "SRDBot works!" {
		t.Errorf("Incorrect message format: %s", body.Text)
	}
}

func TestError(t *testing.T) {
	err := fmt.Errorf("test")
	response := errorResponse(err)

	if response.StatusCode != 500 {
		t.Errorf("Incorrect status code: %d", response.StatusCode)
	}

	if response.Body != "{\"error\": \"test\"}" {
		t.Errorf("Incorrect message format: %s", response.Body)
	}
}
