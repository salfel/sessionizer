package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getProject() (string, bool) {
	projects := findProjects()

	cmd := exec.Command("fzf")
	cmd.Stdin = bytes.NewBufferString(strings.Join(projects, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				fmt.Println("No project found")
				return "", false
			} else if exitErr.ExitCode() == 130 {
				fmt.Println("No project found")
				return "", false
			}

			fmt.Println("Error running fzf:", exitErr)
		}

		fmt.Println("Error running fzf:", err)

		return "", false
	}

	project := strings.Trim(string(output), "\n")

	return project, true
}

func findProjects() []string {
	dir, err := os.ReadDir(DIRPATH)
	if err != nil {
		panic(err)
	}

	var projects []string

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		if !hasGitRepo(fmt.Sprintf("%s/%s", DIRPATH, entry.Name())) {
			continue
		}

		projects = append(projects, entry.Name())
	}

	return projects
}

func hasGitRepo(path string) bool {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return false
	}

	for _, entry := range dir {
		if entry.Name() == ".git" {
			return true
		}
	}

	return false
}
