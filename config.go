package main

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	Windows []Window `toml:"windows"`
}

type Window struct {
	Name string `toml:"name"`
	Path string `toml:"path"`
	Cmd  string `toml:"cmd"`
}

const CONFIG_FILE = "sessionizer.toml"

func LoadConfig(path string) (bool, Config) {
	path = fmt.Sprintf("%s/%s/%s", DIRPATH, PROJECT_NAME, path)

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return false, Config{}
	}

	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		return false, Config{}
	}

	return true, config
}
