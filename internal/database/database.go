package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	sqliteDriver "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "modernc.org/sqlite"
)

// Open opens (or creates) a SQLite database at the given path and runs all
// pending migrations so the schema is up to date.
func Open(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Verify the connection is alive.
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Resolve the schema directory relative to the executable's working directory.
	schemaDir, err := getSchemaDir()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("resolving schema directory: %w", err)
	}

	// Wrap the *sql.DB with the golang-migrate SQLite driver.
	driver, err := sqliteDriver.WithInstance(db, &sqliteDriver.Config{})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("creating migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", schemaDir),
		"sqlite",
		driver,
	)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("creating migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return db, nil
}

// getSchemaDir returns the absolute path to the schema directory.
func getSchemaDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "internal", "database", "schema"), nil
}
