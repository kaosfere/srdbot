package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"content-type": "text/plain"},
		Body:       "SRDBot Works!",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
