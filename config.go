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
const DEFAULT_CONFIG = `
[[windows]]
name = "Editor"
cmd = "nvim"

[[windows]]
name = "Git"
cmd = "lazygit"

[[windows]]
name = "Terminal"
`

func LoadConfig(path string) (Config, bool) {
	path = fmt.Sprintf("%s/%s", path, CONFIG_FILE)

	data, err := os.ReadFile(path)
	if err != nil {
		data = []byte(DEFAULT_CONFIG)
	}

	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		return Config{}, false
	}

	return config, true
}
