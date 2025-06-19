package main

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	SearchPaths []Path        `toml:"search_paths"`
	MaxDepth    int           `toml:"max_depth"`
	Session     SessionConfig `toml:"session"`
}

type SessionConfig struct {
	Windows []Window `toml:"windows"`
	Active  string   `toml:"active"`
}

type Window struct {
	Name string   `toml:"name"`
	Path Path     `toml:"path"`
	Cmd  []string `toml:"cmd"`
}

const GLOBAL_CONFIG = "sessionizer/config.toml"
const CONFIG_FILE = "sessionizer.toml"

const DEFAULT_MAX_DEPTH = 2

func loadDefaultConfig() (Config, bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config dir:", err)
		return Config{}, false
	}

	defaultConfigPath := Path(configDir).Join(GLOBAL_CONFIG)
	config := Config{MaxDepth: DEFAULT_MAX_DEPTH}
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

func loadConfig(path Path) (SessionConfig, bool) {
	path = path.Join(CONFIG_FILE)

	var localConfig SessionConfig
	err := loadConfigFromPath(path, &localConfig)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error loading config file:", err)
		}
		return SessionConfig{}, false
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

	return nil
}
