package main

import "fmt"

func init() {
	InitConfig()
	config, err := GetConfig()
	if err != nil {
		fmt.Println("warn: error in config file.")
	}
	gConfig = config
}
