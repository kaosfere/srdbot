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

	data := &spellData{}
	message, err := getItem("fly", "../../data/spells.json", data)
	assert.Nil(err)
	assert.Equal("in_channel", message.ResponseType)
	assert.Equal(1, len(message.Attachments))

	if len(message.Attachments) > 0 {
		attachment := message.Attachments[0]
		assert.Equal("Fly", attachment.Title)
		fieldsPresent := make([]string, 0)
		for _, field := range attachment.Fields {
			assert.Equal(fieldValues[field.Title], field.Value, field.Title)
			fieldsPresent = append(fieldsPresent, field.Title)
		}

		for field := range fieldValues {
			assert.Contains(fieldsPresent, field)
		}
	}

}

func TestHandleUnkownSpell(t *testing.T) {
	assert := assert.New(t)

	data := &spellData{}
	message, err := getItem("bogus", "../../data/spells.json", data)
	assert.Nil(err)
	assert.Equal("ephemeral", message.ResponseType)
	assert.Equal(0, len(message.Attachments))
	assert.Equal("spell 'bogus' not found", message.Text)
}

func TestHandleMissingSpellsFile(t *testing.T) {
	data := &spellData{}
	_, err := getItem("fly", "bogus", data)
	assert.Error(t, err)
}
