package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dbPath string) {
	var err error

	// Ensure the directory exists
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		position REAL NOT NULL UNIQUE,
		deadline TEXT
	);
	CREATE TABLE IF NOT EXISTS next_actions (
		id TEXT PRIMARY KEY,
		action TEXT NOT NULL,
		project_id TEXT,
		url TEXT,
		size TEXT CHECK(size IS NULL OR size IN ('small', 'medium', 'big')),
		energy TEXT CHECK(energy IS NULL OR energy IN ('high', 'low')),
		created_at DATETIME NOT NULL,
		completed_at DATETIME,
		position REAL NOT NULL UNIQUE,
		FOREIGN KEY(project_id) REFERENCES projects(id)
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}
