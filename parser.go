package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrLogFieldNumberError = errors.New("log field number error")
	ErrLogFieldTypeError   = errors.New("log field type error")
)

type Parser interface {
	Parse(string) (string, error)
}

type TextParser struct{}

func (t *TextParser) Parse(s string) (string, error) {
	return s, nil
}

type JsonParser struct {
	Config *Config
	index  int
}

func (j *JsonParser) Parse(s string) (string, error) {
	arr := strings.Split(s, j.Config.Files[j.index].Separator)
	if len(arr) < len(j.Config.Files[j.index].Fields) {
		return "", ErrLogFieldNumberError
	}

	m := map[string]any{}
	var err error
	for i, field := range j.Config.Files[j.index].Fields {
		switch field.Type {
		case "string":
			m[field.Name] = arr[i]
		case "int":
			m[field.Name], err = strconv.ParseInt(arr[i], 10, 64)
			if err != nil {
				return "", err
			}
		case "float":
			m[field.Name], err = strconv.ParseFloat(arr[i], 64)
			if err != nil {
				return "", err
			}
		default:
			return "", ErrLogFieldTypeError
		}
	}

	buff, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(buff), nil
}

func NewParser(config *Config, index int) Parser {
	if config.Files[index].Separator != "" {
		return &JsonParser{Config: config, index: index}
	}

	return &TextParser{}
}
