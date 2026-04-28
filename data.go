package main

import (
	"fmt"
	"os"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const DATA_DIR = "sessionizer"
const DATABASE_FILE = ".sessionizer.sqlite"

var db *sql.DB

func loadDatabase() bool {
	data_home, ok := os.LookupEnv("XDG_DATA_HOME")
	if !ok {
		home_dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory")
			return false
		}
		data_home = string(Path(home_dir).Join(".local/share"))
	}
	dirPath := Path(data_home).Join(DATA_DIR)
	path := dirPath.Join(DATABASE_FILE)

	err := os.MkdirAll(dirPath.String(), 0755)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	db, err = sql.Open("sqlite3", path.String())
	if err != nil {
		fmt.Println("Error opening database:", err)
		return false
	}

	createTable()

	return true
}

func createTable() bool {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			path TEXT NOT NULL,
			timestamp DATETIME NOT NULL
		);
	`)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func getProjects() []Project {
	rows, err := db.Query(`
		SELECT 
			path, 
			SUM(
				CASE
					WHEN timestamp >= datetime('now', '-7 days') THEN 8
					WHEN timestamp >= datetime('now', '-30 days') THEN 2
					WHEN timestamp >= datetime('now', '-90 days') THEN 1
					ELSE 0
				END
			) as count
		FROM projects 
		WHERE timestamp >= datetime('now', '-30 days') 
		GROUP BY path
		ORDER BY count
	`)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		rows.Scan(&project.Path, &project.Count)
		project.Name = project.Path.String()
		projects = append(projects, project)
	}

	return projects
}

func insertProject(db *sql.DB, project string) bool {
	_, err := db.Exec(`INSERT INTO projects (path, timestamp) VALUES (?, ?)`, project, time.Now())
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
