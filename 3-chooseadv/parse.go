package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

func parse(data []byte) (map[string]StoryArc, error) {
	m := make(map[string]StoryArc)
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("json: %v", err.Error()))
	}
	err = validate(m)
	return m, err
}

func validate(m map[string]StoryArc) error {
	hasIntro := false
	for key, value := range m {
		if key == "intro" {
			hasIntro = true
		}
		for _, o := range value.Options {
			_, ok := m[o.Arc]
			if ok == false {
				return errors.New(fmt.Sprintf("validate: StoryArc %#v has option leading to missing StoryArc %#v", key, o.Arc))
			}
		}
	}
	if !hasIntro {
		return errors.New(fmt.Sprintf("validate: StoryArc \"intro\" missing"))
	}
	return nil
}
