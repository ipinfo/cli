package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/ipinfo/cli/lib"
)

type Config struct {
	GlobalCache bool
}

// init function
func InitConfig() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if !lib.FileExists(path) {
		gConfig.GlobalCache = true
		SetConfig(gConfig)
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
	var setting Config
	err = json.Unmarshal(file, &setting)
	if err != nil {
		return Config{}, err
	}
	return Config{GlobalCache: setting.GlobalCache}, nil
}
