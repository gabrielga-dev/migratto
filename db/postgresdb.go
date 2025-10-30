package db

import (
	"database/sql"
	"fmt"
	"strings"

	DTO "github.com/gabrielga-dev/migratto/dto"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Connect(config DTO.ConfigDTO) (*sql.DB, error) {
	driver := strings.ToLower(config.DatabaseDriver)
	var dsn string

	switch driver {
	case "postgres", "postgresql":

		dsn = fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
			config.DatabaseUsername,
			config.DatabasePassword,
			config.DatabaseName,
			config.DatabaseHost,
			config.DatabasePort,
			config.Sslmode,
		)

	case "mysql", "mariadb":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true",
			config.DatabaseUsername,
			config.DatabasePassword,
			config.DatabaseHost,
			config.DatabasePort,
			config.DatabaseName,
		)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := sql.Open(driverAlias(driver), dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging DB: %w", err)
	}

	return db, nil
}

func driverAlias(name string) string {
	switch name {
	case "postgresql":
		return "postgres"
	case "mariadb":
		return "mysql"
	default:
		return name
	}
}
