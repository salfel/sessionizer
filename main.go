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

	config, ok := LoadConfig(fmt.Sprintf("%s/%s", DIRPATH, project))

	if ok {
		loadSession(project, &config)

		updateData(project)
	}
}
