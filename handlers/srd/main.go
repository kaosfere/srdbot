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

// srdEntry represents basic entry which can be rendered as a slack attachment.
type srdEntry interface {
	asAttachment() slack.Attachment
}

// srdData represents a generic interfacse for loading and finding SRD data.
type srdData interface {
	load(io.Reader) error
	find(string) (srdEntry, error)
}

// commandConfig holds the location of the data source for a given command type
// along with an struct implementing srdData that's prepared to receive the
// unmarshalled data from the data source.
type commandConfig struct {
	sourceFile string
	dataList   srdData
}

// commandConfigs is a map of command configurations keyed by command name.
type commandConfigs map[string]commandConfig

// makeMessage creates a generic slack message without text but with one or
// more attachments.
func makeMessage(info srdEntry) (slack.Msg, error) {
	message := slack.Msg{
		ResponseType: "in_channel",
		Attachments:  []slack.Attachment{info.asAttachment()},
	}

	return message, nil
}

// makeErrorMessage creates an attachment-free slack message with a text field
// meant to hold an error string.
func makeErrorMessage(err error) (slack.Msg, error) {
	message := slack.Msg{
		ResponseType: "ephemeral",
		Text:         fmt.Sprintf("%s", err),
	}

	return message, nil
}

// makeProxyResponse creates a non-error (status code 200) API Gateway response
// containing a JSON-marshaled slack message payload.
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

// makeProxyErrorResponse takes an error object and creates an API Gateway
// response witn status code 500 containeng the text of the error.
func makeProxyErrorResponse(err error) events.APIGatewayProxyResponse {
	response := events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("{\"error\": \"%s\"}", err),
		StatusCode: 500,
	}

	return response
}

// gitItem takes a name to search for, a source file, and a prepared srdData
// implementer, loads data from the source file indo the data list, then
// tries to find the requested data entry, returning it in a slack message
// if found.
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

// handleCommand uses a provided dictionary of command configurations to parse
// and process the provided command, returning a slack message containing the
// requested data if found.
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

// handleRequest fields the request from the API Gateway, builds a command
// configuration, and passes the input data on to handleCommand for
// processing.
func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var message slack.Msg

	commands := commandConfigs{
		"spell":     commandConfig{"data/spells.json", &spellData{}},
		"condition": commandConfig{"data/conditions.json", &conditionData{}},
		"race":      commandConfig{"data/races.json", &raceData{}},
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
