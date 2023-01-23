package bqs

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Host        string
	DebugLevel  string `json:"debug_level"`
	MapFilePath string `json:"map_filepath"`
}

func LoadConf(confPath string) (*Config, error) {
	bytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var c = Config{}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
