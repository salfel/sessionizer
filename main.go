package main

import "fmt"

func main() {
	config, ok := loadDefaultConfig()
	if !ok {
		return
	}

	project, ok := getProject(config.SearchPaths)

	if !ok {
		return
	}

	localConfig, ok := loadConfig(project.Path)
	if ok {
		config.Windows = localConfig.Windows
	}

	if len(config.Windows) == 0 {
		fmt.Println("No windows found in config")
		return
	}

	loadSession(project, &config)

	updateData(project.Path)
}
