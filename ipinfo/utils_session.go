package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func saveToken(tok string) error {
	// create ipinfo config directory.
	cdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	iiCdir := filepath.Join(cdir, "ipinfo")
	if err := os.MkdirAll(iiCdir, 0700); err != nil {
		return err
	}

	// open token file.
	tokFilePath := filepath.Join(iiCdir, "token")
	tokFile, err := os.OpenFile(
		tokFilePath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0660,
	)
	defer tokFile.Close()
	if err != nil {
		return err
	}

	// write token.
	_, err = tokFile.WriteString(tok)
	if err != nil {
		return err
	}

	return nil
}

func deleteToken() error {
	// get token file path.
	cdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	tokFilePath := filepath.Join(cdir, "ipinfo", "token")

	// remove token file.
	if err := os.Remove(tokFilePath); err != nil {
		return err
	}

	return nil
}

func restoreToken() (string, error) {
	// open token file.
	cdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	tokFilePath := filepath.Join(cdir, "ipinfo", "token")
	tokFile, err := os.Open(tokFilePath)
	defer tokFile.Close()
	if err != nil {
		return "", err
	}

	tok, err := ioutil.ReadAll(tokFile)
	if err != nil {
		return "", nil
	}

	return string(tok[:]), nil
}

func fileExists(pathToFile string) bool {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		return false
	}
	return true
}
