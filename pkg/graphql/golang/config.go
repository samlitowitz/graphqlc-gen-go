package golang

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Package   string                `json:"package,omitempty"`
	Output    string                `json:"output,omitempty"`
	ScalarMap map[string]ScalarType `json:"scalarMap,omitempty"`
}

type ScalarType struct {
	Package string `json:"package,omitempty"`
	Type    string `json:"type,omitempty"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
