package main

import "os"

func init() {
	appToken = os.Getenv("APP_TOKEN")
	if appToken == "" {
		panic("APP_TOKEN is not set")
	}
}
