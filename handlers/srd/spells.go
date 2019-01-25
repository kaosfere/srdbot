package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/nlopes/slack"
)

type spellInfo struct {
	Name          string
	Desc          string
	HigherLevel   string `json:"higher_level"`
	Page          string
	Range         string
	Components    string
	Material      string
	Ritual        string
	Duration      string
	Concentration string
	CastingTime   string `json:"casting_time"`
	Level         string
	LevelInt      int `json:"level_int"`
	School        string
	Class         string
	Archetype     string
	Circles       string
}
type spellList []spellInfo

func makeAttachment(spell spellInfo) slack.Attachment {
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
		slack.AttachmentField{
			Title: "Casting Time",
			Value: spell.CastingTime,
			Short: true,
		},
		slack.AttachmentField{
			Title: "Range",
			Value: spell.Range,
			Short: true,
		},
		slack.AttachmentField{
			Title: "Components",
			Value: spell.Components,
			Short: true,
		},
		slack.AttachmentField{
			Title: "Duration",
			Value: spell.Range,
			Short: true,
		},
	}

	if spell.HigherLevel != "" {
		attachment.Fields = append(attachment.Fields,
			slack.AttachmentField{
				Title: "At Higher Levels",
				Value: spell.HigherLevel,
			},
		)
	}

	return attachment
}

func getSpell(name string) (slack.Attachment, error) {
	var spells spellList
	var spellAttachment slack.Attachment

	name = strings.ToLower(name)

	content, err := ioutil.ReadFile("data/spells.json")
	if err != nil {
		return spellAttachment, err
	}

	err = json.Unmarshal(content, &spells)
	if err != nil {
		return spellAttachment, err
	}

	for _, spell := range spells {
		if strings.ToLower(spell.Name) == name {
			spellAttachment = makeAttachment(spell)
			return spellAttachment, nil
		}
	}

	return spellAttachment, nil
}

func handleSpell(name string) (slack.Msg, error) {
	var message slack.Msg

	spellAttachment, err := getSpell(name)
	if err != nil {
		return message, err
	}

	if spellAttachment.Title == "" {
		message = slack.Msg{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf("Spell '%s' not found.", name),
		}
	} else {
		message = slack.Msg{
			ResponseType: "in_channel",
			Attachments:  []slack.Attachment{spellAttachment},
		}
	}

	return message, nil
}
