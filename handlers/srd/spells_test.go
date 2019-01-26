package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleKnownSpell(t *testing.T) {
	fieldValues := map[string]interface{}{
		"Level":            "3rd-level",
		"Class":            "Sorcerer, Warlock, Wizard",
		"Casting Time":     "1 action",
		"Range":            "Touch",
		"Components":       "V, S, M",
		"Duration":         "Up to 10 minutes",
		"At Higher Levels": "When you cast this spell using a spell slot of 4th level or higher, you can target one additional creature for each slot level above 3rd.",
	}

	assert := assert.New(t)

	message, err := handleSpell("fly", "../../data/spells.json")
	assert.Nil(err)
	assert.Equal("in_channel", message.ResponseType)
	assert.Equal(1, len(message.Attachments))

	if len(message.Attachments) > 0 {
		attachment := message.Attachments[0]
		assert.Equal("Fly", attachment.Title)
		for _, field := range attachment.Fields {
			assert.Equal(fieldValues[field.Title], field.Value, field.Title)
		}
	}

}

func TestHandleUnkownSpell(t *testing.T) {
	assert := assert.New(t)

	message, err := handleSpell("bogus", "../../data/spells.json")
	assert.Nil(err)
	assert.Equal("ephemeral", message.ResponseType)
	assert.Equal(0, len(message.Attachments))
	assert.Equal("Spell 'bogus' not found.", message.Text)
}

func TestHandleMissingSpellsFile(t *testing.T) {
	_, err := handleSpell("fly", "bogus")
	assert.Error(t, err)
}
