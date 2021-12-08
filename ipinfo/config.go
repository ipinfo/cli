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
	CacheEnabled bool   `json:"cache_enabled"`
	Token        string `json:"token"`
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

// returns the path to the config file.
func ConfigPath() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confDir, "config.json"), nil
}

// returns the path to the token file.
//
// might be deleted in future release as `token` file is deprecated.
func TokenPath() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confDir, "token"), nil
}

func InitConfig() error {
	configpath, err := ConfigPath()
	if err != nil {
		return err
	}

	// create default config if none yet.
	if !lib.FileExists(configpath) {
		gConfig = NewConfig()

		tokenpath, err := TokenPath()
		if err != nil {
			return err
		}

		// if token exists, migrate it to the config file.
		if lib.FileExists(tokenpath) {
			token, err := ReadTokenFile()
			if err != nil {
				return err
			}
			gConfig.Token = token

			// remove the existing token file
			if err := os.Remove(tokenpath); err != nil {
				return err
			}
		}
	} else {
		config, err := ReadConfig()
		if err != nil {
			return err
		}
		gConfig = config
	}

	if err := SaveConfig(gConfig); err != nil {
		return err
	}

	return nil
}

// returns a new, default config.
func NewConfig() Config {
	return Config{
		CacheEnabled: true,
		Token:        "",
	}
}

// reads `token` file for migration of token to config file.
//
// might be deleted in future release as `token` file is deprecated.
func ReadTokenFile() (string, error) {
	path, err := TokenPath()
	if err != nil {
		return "", err
	}

	token, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// save the config to file.
func SaveConfig(config Config) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

// returns the values of config file.
func ReadConfig() (Config, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
