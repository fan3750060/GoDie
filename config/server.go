package config

import (
	"os"
)

type Http_Config struct {
	AppName string
	Port    string
}

var HttpConfig *Http_Config

func init() {
	loadHttpConfig()
}

func loadHttpConfig() {
	HttpConfig = &Http_Config{
		AppName: os.Getenv("HTTP_APP_NAME"),
		Port:    os.Getenv("HTTP_PORT"),
	}
}
