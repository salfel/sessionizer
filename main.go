package main

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

	loadSession(project, &config)

	updateData(project.Path)
}
