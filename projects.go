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

	projects := make(map[string]Project, 0)
	for _, path := range config.SearchPaths {
		for _, project := range discoverProjects(path, config.MaxDepth) {
			if project.Name == "" {
				fmt.Println("2newline in string", project.Path, project.Count)
			}
			if _, ok := projects[project.Name]; ok && projects[project.Name].Path != project.Path {
				oldProject := projects[project.Name]
				delete(projects, project.Name)

				for _, project := range []Project{oldProject, project} {
					if !strings.Contains(project.Name, "/") {
						segments := project.Path.Segments()
						prefix := segments[len(segments)-2]
						project.Name = fmt.Sprintf("%s/%s", prefix, project.Name)
					}
					projects[project.Name] = project
				}
			} else {
				projects[project.Name] = project
			}
		}
	}

	projectList := make([]Project, 0)
	for _, project := range projects {
		projectList = append(projectList, project)
	}

	sort.Slice(projectList, func(i, j int) bool {
		if projectList[i].Count != projectList[j].Count {
			return projectList[i].Count > projectList[j].Count
		}

		return projectList[i].Name < projectList[j].Name
	})

	projectNames := make([]string, 0)
	for _, project := range projectList {
		projectNames = append(projectNames, project.Name)
		if project.Name == "" {
			fmt.Println("newline in string", project.Path, project.Count)
		}
	}

	if len(projectNames) == 0 {
		fmt.Println("No projects found")
		return Project{}, false
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = bytes.NewBufferString(strings.Join(projectNames, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				fmt.Println("No project found")
				return Project{}, false
			} else if exitErr.ExitCode() == 130 {
				fmt.Println("No project selected")
				return Project{}, false
			}

			fmt.Println("Error running fzf:", exitErr)
		}

		fmt.Println("Error running fzf:", err)

		return Project{}, false
	}

	projectName := strings.Trim(string(output), "\n")

	return projects[projectName], true
}

func findProjects(searchPath Path, data map[Path]int, depth int, maxDepth int) {
	if depth > maxDepth {
		return
	}

	dir, err := os.ReadDir(string(searchPath))
	if err != nil {
		panic(err)
	}

	for _, entry := range dir {
		if entry.IsDir() && entry.Name() == ".git" {
			if _, ok := data[searchPath]; !ok {
				data[searchPath] = 0
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
	data := loadData()

	findProjects(searchPath, data, 0, maxDepth)

	projects := make([]Project, 0)
	for project := range data {
		name := project.Segments()[len(project.Segments())-1]
		if strings.Contains(name, "/") {
			name = strings.Split(name, "/")[1]
		}
		projects = append(projects, Project{Name: name, Path: Path(project), Count: data[project]})
	}

	return projects
}

type ByCount struct {
	counts map[Path]int
	keys   []Path
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
