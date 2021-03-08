package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type App struct {
	TgProdHost  string `yaml:"tgProdHost"`
	TgDevHost   string `yaml:"tgDevHost"`
	AppID       int    `yaml:"appID"`
	AppHash     string `yaml:"appHash"`
	SessionFile string `yaml:"sessionFile"`
	PublicKeys  string `yaml:"publicKeys"`
	PhoneNumber string `yaml:"phoneNumber"`
}

func ReadFromFile(path string) (*App, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	decoder := yaml.NewDecoder(f)
	var app = new(App)
	err = decoder.Decode(&app)
	if err != nil {
		return nil, fmt.Errorf("cannot decode config: %w", err)
	}
	return app, nil
}
