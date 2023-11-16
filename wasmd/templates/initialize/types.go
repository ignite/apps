package initialize

import (
	"gopkg.in/yaml.v2"
)

type Changes struct {
	ChangeItems []ChangeItem `yaml:"changes"`
}
type ChangeItem struct {
	PlaceHolder string `yaml:"placeholder,flow"`
	Text        string `yaml:"text"`
}

type Replaces struct {
	ReplaceItems []ReplaceItem `yaml:"replace"`
}

type ReplaceItem struct {
	Before string `yaml:"before"`
	After  string `yaml:"after"`
}

func parsePlaceHoldersFromYaml(path string) (Changes, error) {
	changes := Changes{}
	source, err := placeholders.ReadFile(path)
	if err != nil {
		return Changes{}, err
	}
	if err := yaml.Unmarshal(source, &changes); err != nil {
		return Changes{}, err
	}
	return changes, nil
}

func parseReplacerFromYaml(path string) (Replaces, error) {
	replaces := Replaces{}

	source, err := placeholders.ReadFile(path)
	if err != nil {
		return Replaces{}, err
	}

	if err := yaml.Unmarshal(source, &replaces); err != nil {
		return Replaces{}, err
	}
	return replaces, nil
}
