package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	err := fmt.Errorf("test")
	response := makeProxyErrorResponse(err)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "{\"error\": \"test\"}", response.Body)
}

func TestBadFormData(t *testing.T) {
	var request events.APIGatewayProxyRequest
	response, err := handleRequest(request)
	assert.Nil(t, err)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "{\"error\": \"Malformed command body\"}", response.Body)
}

func TestUnknownCommandType(t *testing.T) {
	var bodyText map[string]interface{}
	var request events.APIGatewayProxyRequest

	request.Body = "text=bogus"
	response, err := handleRequest(request)
	assert := assert.New(t)
	assert.Nil(err)
	assert.Equal(200, response.StatusCode)

	err = json.Unmarshal([]byte(response.Body), &bodyText)
	if !assert.Nil(err) {
		t.Log("Unmarshal failed.  Unable to continue.")
		t.FailNow()
	}

	assert.Equal("Unknown subcommand: bogus", bodyText["text"])
	assert.Equal("ephemeral", bodyText["response_type"])
	assert.False(bodyText["replace_original"].(bool))
	assert.False(bodyText["delete_original"].(bool))
}

/* func TestSpellCommand(t *testing.T) {
	// This fails because of the JSON file path.
	// Might need to be further paramaterized or something.
	var request events.APIGatewayProxyReques
	request.Body = "text=spell"

	response, err := handleRequest(request)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "", response.Body)
}
*/
