package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleKnownCondition(t *testing.T) {
	assert := assert.New(t)
	config := commandConfigs{"condition": commandConfig{"../../data/conditions.json", &conditionData{}}}
	message, err := handleCommand("condition", "incapacitated", config)

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
	config := commandConfigs{"condition": commandConfig{"../../data/conditions.json", &conditionData{}}}
	message, err := handleCommand("condition", "bogus", config)

	assert.Nil(err)
	assert.Equal("ephemeral", message.ResponseType)
	assert.Equal(0, len(message.Attachments))
	assert.Equal("condition 'bogus' not found", message.Text)
}

func TestHandleMissingConditionFile(t *testing.T) {
	config := commandConfigs{"condition": commandConfig{"bogus", &conditionData{}}}
	_, err := handleCommand("condition", "incapacitated", config)
	assert.Error(t, err)
}
