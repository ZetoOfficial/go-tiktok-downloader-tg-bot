package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	BotToken  string `yaml:"bot_token"`
	DouyinAPI string `yaml:"douyin_api"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
