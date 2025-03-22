package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./gsd.db")
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
