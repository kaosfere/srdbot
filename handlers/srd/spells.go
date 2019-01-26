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
			Value: spell.Duration,
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
