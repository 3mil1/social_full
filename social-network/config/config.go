package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"social-network/pkg/logger"
	"sync"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DbPath string `yaml:"dbPath"`
	} `yaml:"database"`
}

var instance *Config

var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger.InfoLogger.Println("read application configuration")
		instance = &Config{}

		yamlFile, err := ioutil.ReadFile("./config/config.yml")
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

	})
	return instance
}
