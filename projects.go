package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

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

func getProject(searchPaths []Path) (Project, bool) {
	if len(searchPaths) == 0 {
		fmt.Println("No search paths found, please specify at least one")
		return Project{}, false
	}

	projects := make(map[string]Project, 0)
	for _, path := range searchPaths {
		for _, project := range findProjects(path) {
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

	cmd := exec.Command("fzf")
	cmd.Stdin = bytes.NewBufferString(strings.Join(projectNames, "\n"))

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

	projectName := strings.Trim(string(output), "\n")

	return projects[projectName], true
}

func findProjects(searchPath Path) []Project {
	dir, err := os.ReadDir(string(searchPath))
	if err != nil {
		panic(err)
	}

	data := loadData()

	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}

		path := Path(fmt.Sprintf("%s/%s", searchPath, entry.Name()))

		if !hasGitRepo(string(path)) {
			continue
		}

		if _, ok := data[path]; !ok {
			data[path] = 0
		}
	}

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
