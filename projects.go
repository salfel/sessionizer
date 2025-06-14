package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Project struct {
	Name string
	Path string
}

func getProject(searchPaths []string) (Project, bool) {
	if len(searchPaths) == 0 {
		fmt.Println("No search paths found, please specify at least one")
		return Project{}, false
	}

	projects := make([]string, 0)
	for _, path := range searchPaths {
		projects = append(projects, findProjects(path)...)
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = bytes.NewBufferString(strings.Join(projects, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				fmt.Println("No project found")
				return Project{}, false
			} else if exitErr.ExitCode() == 130 {
				fmt.Println("No project found")
				return Project{}, false
			}

			fmt.Println("Error running fzf:", exitErr)
		}

		fmt.Println("Error running fzf:", err)

		return Project{}, false
	}

	project := strings.Trim(string(output), "\n")

	return Project{Name: project, Path: fmt.Sprintf("%s/%s", searchPaths[0], project)}, true
}

func findProjects(searchPath string) []string {
	dir, err := os.ReadDir(searchPath)
	if err != nil {
		panic(err)
	}

	projects := loadData()

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		if !hasGitRepo(fmt.Sprintf("%s/%s", searchPath, entry.Name())) {
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
