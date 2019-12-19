package main

import (
	"io/ioutil"
	"os"

	jsoniter "github.com/json-iterator/go"
)

// JSON instance
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

// DisplayItemConfig display item config struct
type DisplayItemConfig struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// Config struct for config info
type Config struct {
	SupportedSchema      []string                       `json:"supportedSchema"`
	SupportedSearchField map[string][]string            `json:"supportedSearchField"`
	DataPath             map[string]string              `json:"data"`
	DisplayField         map[string][]DisplayItemConfig `json:"display"`
}

// GetConfig stream config from file
func GetConfig() (*Config, error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	result := Config{}
	if err := JSON.Unmarshal(byteValue, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
