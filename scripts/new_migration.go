package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Check if migration name argument is provided
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Migration name argument is required\n")
		os.Exit(1)
	}

	migrationName := os.Args[1]

	// Create timestamp (unix timestamp in milliseconds, first 13 digits)
	timestamp := time.Now().UnixMilli()

	// Create migration file path
	migrationFile := filepath.Join("internal", "sql", "migration", fmt.Sprintf("%d_%s.sql", timestamp, migrationName))

	// Ensure the directory exists
	dir := filepath.Dir(migrationFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Create the migration file
	file, err := os.Create(migrationFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating migration file: %v\n", err)
		os.Exit(1)
	}
	file.Close()

	fmt.Printf("Created migration file: %s\n", migrationFile)
}
