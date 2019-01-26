package main

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/nlopes/slack"
)

type conditionInfo struct {
	Name string
	Desc string
}

func (c conditionInfo) asAttachment() slack.Attachment {
	var attachment slack.Attachment
	attachment.Title = c.Name
	attachment.Text = c.Desc
	return attachment
}

func (c conditionInfo) name() string {
	return c.Name
}

func loadConditions(source io.Reader) ([]conditionInfo, error) {
	var conditions []conditionInfo

	data, err := ioutil.ReadAll(source)
	if err != nil {
		return conditions, err
	}

	err = json.Unmarshal(data, &conditions)
	return conditions, err
}
