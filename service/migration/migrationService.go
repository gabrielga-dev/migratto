package migration_service

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"github.org/gabrielga-dev/migratto/db"
	DTO "github.org/gabrielga-dev/migratto/dto"
	file_service "github.org/gabrielga-dev/migratto/service/file"
)

func Migrate(config DTO.ConfigDTO) error {
	if config.Log {
		fmt.Println("Running migrations...")
	}

	files, err := file_service.GetFilesFromDir(config.MigrationsDir)
	if err != nil {
		fmt.Println("Error getting migration files:", err)
		return err
	}
	err = migrateFiles(files, config)
	if err != nil {
		fmt.Println("\t", err)
		return err
	}
	fmt.Println("Migrations completed successfully.")
	return nil
}

func migrateFiles(files []os.DirEntry, config DTO.ConfigDTO) error {
	databaseConection, err := db.Connect(config)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return err
	}

	err = prepareDatabase(databaseConection)
	if err != nil {
		return err
	}

	appliedMigrations, err := getAppliedMigrations(databaseConection)
	if err != nil {
		return err
	}

	for index, file := range files {
		if len(appliedMigrations) == 0 || index >= len(appliedMigrations) {
			err := migrateFile(file, databaseConection, config)
			if err != nil {
				fmt.Println("Error migrating file:", file.Name())
				return err
			}
		}
	}
	databaseConection.Close()
	return nil
}

func prepareDatabase(databaseConection *sql.DB) error {
	migrattoTableCreationQuery := `
	CREATE TABLE IF NOT EXISTS migratto_migration_history (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`
	_, err := databaseConection.Exec(migrattoTableCreationQuery)
	if err != nil {
		return fmt.Errorf("error creating migratto_migration_history table: %v", err)
	}
	return nil
}

func getAppliedMigrations(databaseConection *sql.DB) ([]string, error) {
	rows, err := databaseConection.Query("SELECT filename FROM migratto_migration_history")
	if err != nil {
		return nil, fmt.Errorf("error querying applied migrations: %v", err)
	}
	defer rows.Close()

	var appliedMigrations []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, fmt.Errorf("error scanning applied migration: %v", err)
		}
		appliedMigrations = append(appliedMigrations, filename)
	}
	return appliedMigrations, nil
}

func migrateFile(file os.DirEntry, databaseConection *sql.DB, config DTO.ConfigDTO) error {
	if config.Log {
		fmt.Println("Migrating file:", file.Name())
	}
	content, err := ioutil.ReadFile(config.MigrationsDir + "/" + file.Name())
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %v", file.Name(), err)
	}
	_, err = databaseConection.Exec(string(content))
	if err != nil {
		return fmt.Errorf("error executing migration file %s: %v", file.Name(), err)
	}
	err = createMigrationHistory(file, databaseConection)
	if err != nil {
		return fmt.Errorf("error creating migration history for file %s: %v", file.Name(), err)
	}
	return nil
}

func createMigrationHistory(file os.DirEntry, databaseConection *sql.DB) error {
	_, err := databaseConection.Exec("INSERT INTO migratto_migration_history (filename) VALUES ($1)", file.Name())
	if err != nil {
		return fmt.Errorf("error recording migration history for file %s: %v", file.Name(), err)
	}
	return nil
}
