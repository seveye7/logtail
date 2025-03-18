package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Name  string `yaml:"name"`
	Files []struct {
		Name      string `yaml:"name"`
		Separator string `yaml:"separator"`
		Fields    []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
		} `yaml:"fields"`
		Topic string `yaml:"topic"`
	} `yaml:"files"`
	Out struct {
		Kafka struct {
			Hosts    []string `yaml:"hosts"`
			Sasl     string   `yaml:"sasl"`
			Username string   `yaml:"username"`
			Password string   `yaml:"password"`
			BatchMax int64    `yaml:"batchMax"`
		} `yaml:"kafka"`
	} `yaml:"out"`
}

func LoadConfig(fileName string) (*Config, error) {
	temp := &Config{}
	//
	buff, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buff, temp)
	if err != nil {
		return nil, err
	}

	return temp, nil
}
