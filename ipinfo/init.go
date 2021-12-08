package main

import "fmt"

func init() {
	if err := InitConfig(); err != nil {
		fmt.Println("warn: error in creating config file.")
	}
}
