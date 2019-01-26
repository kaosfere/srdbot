package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
)

type raceInfo struct {
	Name      string
	ASI       string `json:"asi-desc"`
	Age       string
	Alignment string
	Size      string
	Speed     string `json:"speed-desc"`
	Languages string
	Traits    string
	Subtypes  []raceSubtype
}

type raceSubtype struct {
	Name   string
	Desc   string
	ASI    string `json:"asi-desc"`
	Traits string
}

type raceData struct {
	races []raceInfo
}

func (r *raceData) load(source io.Reader) error {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &r.races)
	return err
}

func (r raceData) find(name string) (srdEntry, error) {
	var race raceInfo
	for _, race := range r.races {
		if strings.ToLower(race.Name) == strings.ToLower(name) {
			return race, nil
		}
	}
	return race, fmt.Errorf("race '%s' not found", name)
}

func (race raceInfo) asAttachment() slack.Attachment {
	var attachment slack.Attachment

	headerRegex := regexp.MustCompile("\\*\\*_(.*)\\._\\*\\*")
	//boldRegex := regexp.MustCompile("\\*\\*(.*)\\*\\*")
	asi := headerRegex.ReplaceAllString(race.ASI, "*$1:*")
	age := headerRegex.ReplaceAllString(race.Age, "*$1:*")
	alignment := headerRegex.ReplaceAllString(race.Alignment, "*$1:*")
	size := headerRegex.ReplaceAllString(race.Size, "*$1:*")
	speed := headerRegex.ReplaceAllString(race.Speed, "*$1:*")
	languages := headerRegex.ReplaceAllString(race.Languages, "*$1:*")

	details := []string{asi, age, alignment, size, speed, languages}

	if race.Traits != "" {
		traits := headerRegex.ReplaceAllString(race.Traits, "*$1:*")
		details = append(details, traits)
	}

	if len(race.Subtypes) > 0 {
		details = append(details, "\n*Subtypes*")
		for _, subtype := range race.Subtypes {
			// TODO: Fill this out with subtype data formatted... somehow
			details = append(details, fmt.Sprintf("* _%s_", subtype.Name))
		}
	}
	attachment.Title = race.Name
	attachment.Text = strings.Join(details, "\n")

	return attachment
}
