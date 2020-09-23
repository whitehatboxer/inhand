package main

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type DispatcherConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Interval time.Duration `yaml:"interval"`
}

type CrawlerConfig struct {
	 Goroutines int `yaml:"goroutines"`
}

type Config struct {
	DispatcherConfig `yaml:"dispatcher"`
	CrawlerConfig     `yaml:"crawler"`
}

func ConfigInit(configPath string, conf *Config) error {
	content, readErr := ioutil.ReadFile(configPath)
	if readErr != nil {
		return readErr
	}

	unmarshalErr := yaml.UnmarshalStrict(content, conf)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}
