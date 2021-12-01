package main

import "fmt"

func init() {
	err := InitConfig()
	if err != nil {
		fmt.Println("warn: error in creating config file.")
	}
	config, err := GetConfig()
	if err != nil {
		fmt.Println("warn: error in config file.")
	}
	gConfig = config
}
