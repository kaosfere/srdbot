package main

import (
	"encoding/json"
	"fmt"

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
	message := slack.Msg{
		ResponseType: "ephemeral",
		Text:         "SRDBot works!",
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
