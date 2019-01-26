package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nlopes/slack"
)

type conditionInfo struct {
	Name string
	Desc string
}

type conditionList []conditionInfo

func makeConditionAttachment(condition conditionInfo) slack.Attachment {
	var attachment slack.Attachment
	attachment.Title = condition.Name
	attachment.Text = condition.Desc
	return attachment
}

func getCondition(name string, sourceFile string) (slack.Attachment, error) {
	var conditions conditionList
	var conditionAttachment slack.Attachment

	name = strings.ToLower(name)

	content, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return conditionAttachment, err
	}

	err = json.Unmarshal(content, &conditions)
	if err != nil {
		return conditionAttachment, err
	}

	for _, condition := range conditions {
		if strings.ToLower(condition.Name) == name {
			conditionAttachment = makeConditionAttachment(condition)
			return conditionAttachment, nil
		}
	}

	return conditionAttachment, nil
}

func handleCondition(name string, sourceFile string) (slack.Msg, error) {
	var message slack.Msg

	conditionAttachment, err := getCondition(name, sourceFile)
	if err != nil {
		return message, err
	}

	if conditionAttachment.Title == "" {
		message = slack.Msg{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Condition '%s' not found.", name),
		}
	} else {
		message = slack.Msg{
			ResponseType: "in_channel",
			Attachments:  []slack.Attachment{conditionAttachment},
		}
	}

	return message, nil
}
