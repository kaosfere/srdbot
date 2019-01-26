package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/nlopes/slack"
)

type srdEntry interface {
	asAttachment() slack.Attachment
}

type srdData interface {
	load(io.Reader) error
	find(string) (srdEntry, error)
}

type commandConfig struct {
	sourceFile string
	dataList   srdData
}

type commandConfigs map[string]commandConfig

func makeMessage(info srdEntry) (slack.Msg, error) {
	message := slack.Msg{
		ResponseType: "in_channel",
		Attachments:  []slack.Attachment{info.asAttachment()},
	}

	return message, nil
}

func makeErrorMessage(err error) (slack.Msg, error) {
	message := slack.Msg{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("%s", err),
	}

	return message, nil
}

func makeProxyResponse(message slack.Msg, err error) events.APIGatewayProxyResponse {
	if err != nil {
		return makeProxyErrorResponse(err)
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return makeProxyErrorResponse(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(messageJSON),
		StatusCode: 200,
	}
}

func makeProxyErrorResponse(err error) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", err),
		StatusCode: 500,
	}

	return response
}

func getItem(name string, sourceFile string, dataList srdData) (slack.Msg, error) {
	var message slack.Msg
	source, err := os.Open(sourceFile)
	defer source.Close()
	if err != nil {
		return message, err
	}

	err = dataList.load(source)
	if err != nil {
		return message, err
	}

	item, err := dataList.find(name)
	if err != nil {
		return makeErrorMessage(err)
	}

	return makeMessage(item)
}

func handleCommand(command string, args string, configs commandConfigs) (slack.Msg, error) {
	var err error
	var message slack.Msg

	if config, ok := configs[command]; ok {
		message, err = getItem(args, config.sourceFile, config.dataList)
	} else {
		message = slack.Msg{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Unknown subcommand: %s", command),
		}
		err = nil
	}

	return message, err
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var message slack.Msg

	commands := commandConfigs{
		"spell":     commandConfig{"data/spells.json", &spellData{}},
		"condition": commandConfig{"data/conditions.json", &conditionData{}},
	}

	// Mock up an http request response so we can use ParseForm() to process the body
	httpRequest, err := http.NewRequest("POST", "", strings.NewReader(request.Body))
	if err != nil {
		return makeProxyResponse(message, err), nil
	}

	httpRequest.Header.Set("content-type", "application/x-www-form-urlencoded")

	err = httpRequest.ParseForm()
	if err != nil {
		return makeProxyResponse(message, err), nil
	}

	if len(httpRequest.PostForm["text"]) == 0 {
		return makeProxyResponse(message, fmt.Errorf("Malformed command body")), nil
	}
	requestBody := httpRequest.PostForm["text"][0]
	requestParts := strings.Split(requestBody, " ")

	command := requestParts[0]
	args := strings.Join(requestParts[1:], " ")
	message, err = handleCommand(command, args, commands)

	return makeProxyResponse(message, err), nil
}

func main() {
	lambda.Start(handleRequest)
}
