package common

import "os"

func GetConfig() (*Config) {
	config := Config {
		Token: os.Getenv("TOKEN"),
		SinceId: os.Getenv("SINCE_ID"),
		GroupId: os.Getenv("GROUP_ID"),
	}

	return &config
}