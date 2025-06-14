package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// TODO: Store the data in a sqlite db and add timestamps to it for better filtering

const DATA_FILE = ".sessionizer.toml"

func storeData(projects map[string]int) {
	data, err := toml.Marshal(projects)
	if err != nil {
		fmt.Println("Error marshalling config file:", err)
		return
	}

	err = os.WriteFile(getDataPath(), data, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}
}

func loadData() map[string]int {
	data, err := os.ReadFile(getDataPath())
	if err != nil {
		return map[string]int{}
	}

	var projects map[string]int
	err = toml.Unmarshal(data, &projects)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		return map[string]int{}
	}

	return projects
}

func getDataPath() string {
	homeDir := getHomeDir()

	return fmt.Sprintf("%s/%s", homeDir, DATA_FILE)
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting home directory: %s", err))
	}

	return homeDir
}
