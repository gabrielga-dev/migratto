package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	DTO "github.org/gabrielga-dev/migratto/dto"
)

func Connect(config DTO.ConfigDTO) (*sql.DB, error) {
	connection := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s sslmode=%s",
		config.DatabaseUsername,
		config.DatabaseName,
		config.DatabasePassword,
		config.DatabaseHost,
		config.Sslmode,
	)
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	return db, nil
}
