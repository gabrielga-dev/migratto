package migration_service

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gabrielga-dev/migratto/db"
	DTO "github.com/gabrielga-dev/migratto/dto"
	migration_model "github.com/gabrielga-dev/migratto/model/migration"
	migration_collection_model "github.com/gabrielga-dev/migratto/model/migration/collection"
	file_service "github.com/gabrielga-dev/migratto/service/file"
)

func Migrate(config DTO.ConfigDTO) error {
	if config.Log {
		fmt.Println("Running migrations...")
	}

	err := validateDriver(config.DatabaseDriver)
	if err != nil {
		fmt.Println("Error validating database driver:", err)
		return err
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

func validateDriver(driver string) error {
	supportedDrivers := []string{"postgres", "postgresql", "mysql", "mariadb"}
	for _, d := range supportedDrivers {
		if driver == d {
			return nil
		}
	}
	return fmt.Errorf("%s is an unsupported database driver", driver)
}

func migrateFiles(files []os.DirEntry, config DTO.ConfigDTO) error {
	databaseConection, err := db.Connect(config)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return err
	}

	err = prepareDatabase(databaseConection, config)
	if err != nil {
		return err
	}

	appliedMigrations, err := getAppliedMigrations(databaseConection)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := migrateFile(file, databaseConection, config, appliedMigrations)
		if err != nil {
			fmt.Println("Error migrating file:", file.Name())
			return err
		}
	}
	databaseConection.Close()
	return nil
}

func prepareDatabase(databaseConection *sql.DB, config DTO.ConfigDTO) error {
	migrattoTableCreationQuery := getMigrattoTableCreationQuery(config.DatabaseDriver)
	_, err := databaseConection.Exec(migrattoTableCreationQuery)
	if err != nil {
		return fmt.Errorf("error creating migratto_migration_history table: %v", err)
	}
	return nil
}

func getMigrattoTableCreationQuery(driver string) string {
	switch driver {
	case "postgres", "postgresql":
		return `
		CREATE TABLE IF NOT EXISTS migratto_migration_history (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL,
			checksum VARCHAR NOT NULL,
			tag VARCHAR NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		`
	case "mysql", "mariadb":
		return `
		CREATE TABLE IF NOT EXISTS migratto_migration_history (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			filename VARCHAR(255) NOT NULL,
			checksum VARCHAR(255) NOT NULL,
			tag VARCHAR(100) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	default:
		return ""
	}
}

func getAppliedMigrations(databaseConection *sql.DB) (migration_collection_model.MigrationModelCollection, error) {
	rows, err := databaseConection.Query("SELECT filename, checksum, tag FROM migratto_migration_history")
	if err != nil {
		return migration_collection_model.MigrationModelCollection{}, fmt.Errorf("error querying applied migrations: %v", err)
	}
	defer rows.Close()

	var appliedMigrations []migration_model.MigrationModel
	for rows.Next() {
		var filename, checksum, tag string

		if err := rows.Scan(&filename, &checksum, &tag); err != nil {
			return migration_collection_model.MigrationModelCollection{}, fmt.Errorf("error scanning applied migration: %v", err)
		}
		migrationModel := migration_model.MigrationModel{
			Filename: filename,
			Checksum: checksum,
			Tag:      tag,
		}

		appliedMigrations = append(appliedMigrations, migrationModel)
	}
	return migration_collection_model.MigrationModelCollection{Migrations: appliedMigrations}, nil
}

func migrateFile(
	file os.DirEntry,
	databaseConection *sql.DB,
	config DTO.ConfigDTO,
	appliedMigrations migration_collection_model.MigrationModelCollection,
) error {

	if config.Log {
		fmt.Println("Migrating file:", file.Name())
	}

	fileTag := file_service.GetFileTag(file.Name())
	appliedMigration, err := appliedMigrations.GetMigrationByTag(fileTag)
	if err == nil {
		// Migration with this tag has been applied
		if !appliedMigration.IsEqual(file, config.MigrationsDir) {
			return fmt.Errorf("migration conflict for tag %s in file %s", fileTag, file.Name())
		} else {
			return nil
		}
	}

	content, err := ioutil.ReadFile(config.MigrationsDir + "/" + file.Name())
	if err != nil {
		return fmt.Errorf("error reading migration file %s: %v", file.Name(), err)
	}

	_, err = databaseConection.Exec(string(content))
	if err != nil {
		return fmt.Errorf("error executing migration file %s: %v", file.Name(), err)
	}

	err = createMigrationHistory(file, databaseConection, config)
	if err != nil {
		return fmt.Errorf("error creating migration history for file %s: %v", file.Name(), err)
	}

	return nil
}

func createMigrationHistory(file os.DirEntry, databaseConection *sql.DB, config DTO.ConfigDTO) error {
	checksum, err := file_service.GetChecksum(config.MigrationsDir + "/" + file.Name())
	if err != nil {
		return fmt.Errorf("error getting checksum for file %s: %v", file.Name(), err)
	}

	tag := file_service.GetFileTag(file.Name())

	query := fmt.Sprintf(
		"INSERT INTO migratto_migration_history (filename, checksum, tag) VALUES ('%s', '%s', '%s')",
		file.Name(), checksum, tag,
	)
	_, err = databaseConection.Exec(query)
	if err != nil {
		return fmt.Errorf("error recording migration history for file %s: %v", file.Name(), err)
	}
	return nil
}
