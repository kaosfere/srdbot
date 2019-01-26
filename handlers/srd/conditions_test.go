package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleKnownCondition(t *testing.T) {
	assert := assert.New(t)

	data := &conditionData{}
	message, err := getItem("incapacitated", "../../data/conditions.json", data)
	assert.Nil(err)
	assert.Equal("in_channel", message.ResponseType)
	assert.Equal(1, len(message.Attachments))

	if len(message.Attachments) > 0 {
		attachment := message.Attachments[0]
		assert.Equal("Incapacitated", attachment.Title)
		assert.Equal("* An incapacitated creature can’t take actions or reactions.", attachment.Text)
	}

}

func TestHandleUnkownCondition(t *testing.T) {
	assert := assert.New(t)

	data := &conditionData{}
	message, err := getItem("bogus", "../../data/conditions.json", data)
	assert.Nil(err)
	assert.Equal("ephemeral", message.ResponseType)
	assert.Equal(0, len(message.Attachments))
	assert.Equal("condition 'bogus' not found", message.Text)
}

func TestHandleMissingConditionFile(t *testing.T) {
	data := &conditionData{}
	_, err := getItem("incapacitated", "bogus", data)
	assert.Error(t, err)
}
