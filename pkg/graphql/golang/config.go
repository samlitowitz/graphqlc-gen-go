package golang

import (
	"encoding/json"
	cfgPkg "github.com/samlitowitz/graphqlc-gen-go/pkg/graphql/golang/config"
	"io/ioutil"
)

type Config struct {
	Package      string                       `json:"package,omitempty"`
	OutputPrefix string                       `json:"output_prefix,omitempty"`
	ScalarMap    map[string]cfgPkg.ScalarType `json:"scalarMap,omitempty"`
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
