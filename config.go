package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Twitter struct {
		ConsumerKey    string
		ConsumerSecret string
		Token          string
		TokenSecret    string
	}
	ImKayacCom struct {
		User     string
		Password string
		Secret   string
	}
	Keywords          []string
	FilterScreenNames []string
	Events            []string
	Agent             []struct {
		Name   string
		Params map[string]string
	}
}

//TODO: build *Config from enviroment value
func NewConfig(filePath string) (*Config, error) {
	byt, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(byt, &config); err != nil {
		return nil, err
	}

	return config, nil
}
