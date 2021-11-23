package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func getTokenFilePath() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confDir, "token"), nil
}

func saveToken(tok string) error {
	tokFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	// open/create if necessary.
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
	tokFilePath, err := getTokenFilePath()
	if err != nil {
		return err
	}

	// remove token file.
	if err := os.Remove(tokFilePath); err != nil {
		return err
	}

	return nil
}

func restoreToken() (string, error) {
	tokFilePath, err := getTokenFilePath()
	if err != nil {
		return "", err
	}

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
