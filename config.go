package main

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
)

// yaml library uses json tags
type ServerConfig struct {
	StripeSecretKey     string `json:"stripeSecretKey"`
	DbName              string `json:"dbName"`
	EmailSenderAddress  string `json:"emailSenderAddress"`
	EmailSenderPassword string `json:"emailSenderPassword"`
	EmailSenderName     string `json:"emailSenderName"`
	AdminAddress        string `json:"adminAddress"`
}

func (config *ServerConfig) load(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, config)
	return err
}

func readConfig(path string) (interface{}, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	var json interface{}
	err = yaml.Unmarshal(yamlFile, &json)
	if err != nil {
		return "", err
	}

	return json, nil
}
