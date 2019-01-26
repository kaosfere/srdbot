package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/nlopes/slack"
)

type conditionInfo struct {
	Name string
	Desc string
}

type conditionData struct {
	conditions []conditionInfo
}

func (c *conditionData) load(source io.Reader) error {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c.conditions)
	return err
}

func (c conditionData) find(name string) (srdEntry, error) {
	var condition conditionInfo
	for _, condition := range c.conditions {
		if strings.ToLower(condition.Name) == strings.ToLower(name) {
			return condition, nil
		}
	}
	return condition, fmt.Errorf("condition '%s' not found", name)
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
