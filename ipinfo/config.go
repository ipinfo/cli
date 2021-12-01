package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ipinfo/cli/lib"
)

// global config.
var gConfig Config

type Config struct {
	Cache bool
}

// gets the global config directory, creating it if necessary.
func getConfigDir() (string, error) {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	confDir := filepath.Join(cdir, "ipinfo")
	if err := os.MkdirAll(confDir, 0700); err != nil {
		return "", err
	}

	return confDir, nil
}

// init function
func InitConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if !lib.FileExists(path) {
		err := SetConfig(NewConfig())
		if err != nil {
			return err
		}
	}
	return nil
}

// returns the path to config file.
func ConfigPath() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(confDir, "config.json"), nil
}

// returns config file with default settings.
func NewConfig() Config {
	return Config{
		Cache: true,
	}
}

// set the values of config.
func SetConfig(config Config) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	jsonFile, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, jsonFile, 0644)
	if err != nil {
		return err
	}
	return nil
}

// get the values of config.
func GetConfig() (Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
