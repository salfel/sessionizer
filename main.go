package main

const DIRPATH = "/home/felix/Projects"
const PROJECT_NAME = "fenix"

func main() {
	exists, config := LoadConfig(CONFIG_FILE)

	if !exists {
		loadSession(PROJECT_NAME, nil)
	} else {
		loadSession(PROJECT_NAME, &config)
	}
}
