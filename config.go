package main

import (
	"fmt"
	"os"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	SearchPaths []Path   `toml:"search_paths"`
	Windows     []Window `toml:"windows"`
}

type WindowConfig struct {
	Windows []Window `toml:"windows"`
}

type Window struct {
	Name string   `toml:"name"`
	Path string   `toml:"path"`
	Cmd  []string `toml:"cmd"`
}

const GLOBAL_CONFIG = "sessionizer/config.toml"
const CONFIG_FILE = "sessionizer.toml"

func loadDefaultConfig() (Config, bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config dir:", err)
		return Config{}, false
	}

	defaultConfigPath := Path(configDir).Join(GLOBAL_CONFIG)
	var config Config
	err = loadConfigFromPath(defaultConfigPath, &config)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("No config file at %s found\n", defaultConfigPath)
			return config, false
		} else {
			fmt.Println("Error loading config file:", err)
			return config, false
		}
	}

	return config, true
}

func loadConfig(path Path) (WindowConfig, bool) {
	path = path.Join(CONFIG_FILE)

	var localConfig WindowConfig
	err := loadConfigFromPath(path, &localConfig)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error loading config file:", err)
		}
		return WindowConfig{}, false
	}

	return localConfig, true
}

func loadConfigFromPath[T any](path Path, config *T) error {
	data, err := os.ReadFile(string(path))
	if err != nil {
		return err
	}

	err = toml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	var windows []Window
	switch c := any(config).(type) {
	case *WindowConfig:
		windows = c.Windows
	case *Config:
		windows = c.Windows
	}

	for i, window := range windows {
		for j, cmd := range window.Cmd {
			if strings.Contains(cmd, " ") {
				windows[i].Cmd[j] = fmt.Sprintf("'%s'", cmd)
			}
		}
	}

	return nil
}
