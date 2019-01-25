package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nlopes/slack"
)

func errorResponse(err error) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", err),
		StatusCode: 500,
	}

	return response
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var message slack.Msg

	// Mock up an http request response so we can use ParseForm() to process the body
	httpRequest, err := http.NewRequest("POST", "", strings.NewReader(request.Body))
	if err != nil {
		return errorResponse(err), nil

	}
	httpRequest.Header.Set("content-type", "application/x-www-form-urlencoded")

	err = httpRequest.ParseForm()
	if err != nil {
		fmt.Printf("ERR: %s\n", err)
		return errorResponse(err), nil
	}

	requestBody := httpRequest.Form["text"][0]
	requestParts := strings.Split(requestBody, " ")

	commandType := requestParts[0]

	switch commandType {
	case "spell":
		message, err = handleSpell(strings.Join(requestParts[1:], " "))
	default:
		message = slack.Msg{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Unknown subcommand: %s\nCurrently supported: spell", requestParts[0]),
		}
		err = nil
	}

	if err != nil {
		return errorResponse(err), nil
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return errorResponse(err), nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string(messageJSON),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
