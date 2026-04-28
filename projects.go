package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strings"
)

var excluded_directories = []string{"node_modules", "vendor", "target", "build", "dist", "out"}

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) Join(path string) Path {
	return Path(fmt.Sprintf("%s/%s", p.String(), path))
}

func (p Path) Segments() []string {
	return strings.Split(p.String(), "/")
}

type Project struct {
	Name  string
	Path  Path
	Count int
}

func getProject(config Config) (Project, bool) {
	if len(config.SearchPaths) == 0 {
		fmt.Println("No search paths found, please specify at least one")
		return Project{}, false
	}

	projects := discoverProjects(config.SearchPaths[0], config.MaxDepth)

	sort.Slice(projects, func(i, j int) bool {
		if projects[i].Count != projects[j].Count {
			return projects[i].Count > projects[j].Count
		}

		return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
	})

	projectNames := make([]string, 0, len(projects))
	for _, project := range projects {
		projectNames = append(projectNames, project.Name)
		if project.Name == "" {
			fmt.Println("newline in string", project.Path, project.Count)
		}
	}

	if len(projectNames) == 0 {
		fmt.Println("No projects found")
		return Project{}, false
	}

	projectName, ok := launchFzf(projectNames)
	if !ok {
		return Project{}, false
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project, true
		}
	}

	return Project{}, false
}

func launchFzf(names []string) (string, bool) {
	cmd := exec.Command("fzf")
	cmd.Stdin = bytes.NewBufferString(strings.Join(names, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				fmt.Println("No project found")
				return "", false
			} else if exitErr.ExitCode() == 130 {
				fmt.Println("No project selected")
				return "", false
			}

			fmt.Println("Error running fzf:", exitErr)
		}

		fmt.Println("Error running fzf:", err)

		return "", false
	}

	return strings.Trim(string(output), "\n"), true
}

func findProjects(searchPath Path, data *[]Project, depth int, maxDepth int) {
	if depth > maxDepth {
		return
	}

	dir, err := os.ReadDir(string(searchPath))
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {
		if entry.IsDir() && entry.Name() == ".git" {
			found := false
			for _, project := range *data {
				if project.Path == searchPath {
					found = true
					break
				}
			}
			if !found {
				*data = append(*data, Project{Name: searchPath.String(), Path: searchPath, Count: 0})
			}

			return
		}
	}

	for _, entry := range dir {
		if !entry.IsDir() || slices.Contains(excluded_directories, entry.Name()) {
			continue
		}

		path := Path(searchPath).Join(entry.Name())

		findProjects(path, data, depth+1, maxDepth)
	}
}

func discoverProjects(searchPath Path, maxDepth int) []Project {
	data := getProjects()

	findProjects(searchPath, &data, 0, maxDepth)

	for i, project := range data {
		name := strings.Replace(project.Path.String(), searchPath.String()+"/", "", 1)
		data[i].Name = name
	}

	return data
}
