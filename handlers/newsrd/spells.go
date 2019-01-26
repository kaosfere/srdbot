package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/nlopes/slack"
)

type spellInfo struct {
	Name  string
	Desc  string
	Level string
	Class string
}

type spellData struct {
	spells []spellInfo
}

func (s *spellData) load(source io.Reader) error {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &s.spells)
	return err
}

func (s spellData) find(name string) (srdEntry, error) {
	var spell spellInfo
	for _, spell := range s.spells {
		if strings.ToLower(spell.Name) == strings.ToLower(name) {
			return spell, nil
		}
	}
	return spell, fmt.Errorf("spell '%s' not found", name)
}

type spellList []spellInfo

func (spell spellInfo) asAttachment() slack.Attachment {
	var attachment slack.Attachment
	attachment.Title = spell.Name
	attachment.Text = spell.Desc

	attachment.Fields = []slack.AttachmentField{
		slack.AttachmentField{
			Title: "Level",
			Value: spell.Level,
			Short: true,
		},
		slack.AttachmentField{
			Title: "Class",
			Value: spell.Class,
			Short: true,
		},
	}
	return attachment
}

func (spell spellInfo) name() string {
	return spell.Name
}

func loadSpells(source io.Reader) ([]spellInfo, error) {
	var spells []spellInfo

	data, err := ioutil.ReadAll(source)
	if err != nil {
		return spells, err
	}

	err = json.Unmarshal(data, &spells)
	return spells, err
}
