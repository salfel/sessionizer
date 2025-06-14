package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
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

	data := loadData()
	if _, ok := data[project]; !ok {
		data[project] = 1
	} else {
		data[project]++
	}
	storeData(data)

	return project, true
}

func findProjects() []string {
	dir, err := os.ReadDir(DIRPATH)
	if err != nil {
		panic(err)
	}

	projects := loadData()

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		if !hasGitRepo(fmt.Sprintf("%s/%s", DIRPATH, entry.Name())) {
			continue
		}

		if _, ok := projects[entry.Name()]; !ok {
			projects[entry.Name()] = 0
		}
	}

	keys := make([]string, 0)
	for project := range projects {
		keys = append(keys, project)
	}

	sortedProjects := ByCount{counts: projects, keys: keys}
	sort.Sort(sortedProjects)

	return sortedProjects.keys
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

type ByCount struct {
	counts map[string]int
	keys   []string
}

func (c ByCount) Len() int {
	return len(c.keys)
}
func (c ByCount) Less(i, j int) bool {
	return c.counts[c.keys[i]] > c.counts[c.keys[j]]
}
func (c ByCount) Swap(i, j int) {
	c.keys[i], c.keys[j] = c.keys[j], c.keys[i]
}
