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
	requestText, err := json.MarshalIndent(request, "", "    ")
	if err != nil {
		return errorResponse(err), nil
	}

	// Mock up an http request response so we can use ParseForm() to process the body
	httpRequest, err := http.NewRequest("POST", "", strings.NewReader(request.Body))
	if err != nil {
		return errorResponse(err), nil

	}
	httpRequest.Header.Set("content-type", "application/x-www-form-urlencoded")

	err = httpRequest.ParseForm()
	if err != nil {
		return errorResponse(err), nil
	}

	message := slack.Msg{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("Request is: ```%s1```", string(requestText)),
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
