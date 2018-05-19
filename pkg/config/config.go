package config

import (
	"io/ioutil"

	log "github.com/jaredallard/mm/pkg/logger"
	"gopkg.in/yaml.v2"
)

// ConfigurationFile is mm's config file
type ConfigurationFile struct {
	Telegram struct {
		Token     string `yaml:"token"`
		ChannelID string `yaml:"channelID"`
	}
	Sheet struct {
		ID string `yaml:"ID"`
	}
	Google struct {
		APIKey string `yaml:"apiKey"`
	}
}

// Load the configuration file.
func Load(path string) ConfigurationFile {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to load configuration file at", path)
	}

	config := ConfigurationFile{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal("Failed to unmarshal configuration file:", err.Error())
	}

	return config
}
