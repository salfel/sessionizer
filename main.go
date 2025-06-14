package main

import (
	"fmt"
)

const DIRPATH = "/home/felix/Projects"

func main() {
	project, ok := getProject()

	if !ok {
		return
	}

	config, ok := loadConfig(fmt.Sprintf("%s/%s", DIRPATH, project))

	if ok {
		loadSession(project, &config)

		updateData(project)
	}
}
