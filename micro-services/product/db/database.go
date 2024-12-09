package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase(databasePath string) *sql.DB {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	fmt.Println("Database initialized successfully")
	return db
}

func RunMigrations(db *sql.DB, migrationFile string) {
	// Get the directory of the currently running binary
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	baseDir := filepath.Dir(exePath)

	// Combine the baseDir with the relative migrationFile path
	fullPath := filepath.Join(baseDir, migrationFile)

	log.Printf("Running migrations from: %s", fullPath)

	// Open the migration file
	migrationContent, err := os.ReadFile(fullPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	// Execute the migration
	_, err = db.Exec(string(migrationContent))
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	log.Println("Migration executed successfully")
}
