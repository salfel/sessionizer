package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// TODO: Store the data in a sqlite db and add timestamps to it for better filtering

const DATA_FILE = ".sessionizer.toml"

func storeData(projects map[Path]int) {
	data, err := toml.Marshal(projects)
	if err != nil {
		fmt.Println("Error marshalling config file:", err)
		return
	}

	err = os.WriteFile(string(getDataPath()), data, 0644)
	if err != nil {
		fmt.Println("Error writing config file:", err)
		return
	}
}

func loadData() map[Path]int {
	data, err := os.ReadFile(string(getDataPath()))
	if err != nil {
		return map[Path]int{}
	}

	var projects map[Path]int
	err = toml.Unmarshal(data, &projects)
	if err != nil {
		fmt.Println("Error unmarshalling config file:", err)
		return map[Path]int{}
	}

	return projects
}

func updateData(projectPath Path) {
	data := loadData()
	if _, ok := data[projectPath]; !ok {
		data[projectPath] = 1
	} else {
		data[projectPath]++
	}
	storeData(data)
}

func getDataPath() Path {
	homeDir := getHomeDir()

	return Path(homeDir).Join(DATA_FILE)
}

func getHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Error getting home directory: %s", err))
	}

	return homeDir
}
