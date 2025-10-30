package main

import (
	"fmt"

	DTO "github.com/gabrielga-dev/migratto/dto"
	migration_service "github.com/gabrielga-dev/migratto/service/migration"
)

func main() {
	fmt.Println("Welcome to Migratto!")
	migration_service.Migrate(DTO.ConfigDTO{
		DatabaseDriver:   "postgres",
		DatabaseHost:     "localhost",
		DatabasePort:     5432,
		DatabaseName:     "migratto",
		DatabaseUsername: "user",
		DatabasePassword: "pass",
		Sslmode:          "disable",
		MigrationsDir:    "./migrations",
		Log:              true,
	})
}
