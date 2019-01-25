package main

import (
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var response events.APIGatewayProxyResponse
	content, err := ioutil.ReadFile("data/spells.json")
	if err != nil {
		response = events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("{\"error\": \"%s\"}", err),
			StatusCode: 500,
		}
	} else {
		response = events.APIGatewayProxyResponse{
			//Headers:    map[string]string{"content-type": "text/plain"},
			Body:       string(content),
			StatusCode: 500,
		}
	}

	return response, nil
}

func main() {
	lambda.Start(handleRequest)
}
