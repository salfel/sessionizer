package main

import "fmt"

func main() {
	config, ok := loadDefaultConfig()
	if !ok {
		return
	}

	project, ok := getProject(config)

	if !ok {
		return
	}

	localConfig, ok := loadConfig(project.Path)
	if ok {
		config.Session = localConfig
	}

	if len(config.Session.Windows) == 0 {
		fmt.Println("No windows found in config")
		return
	}

	if err := loadSession(project, &config); err != nil {
		return
	}

	updateData(project.Path)
}
