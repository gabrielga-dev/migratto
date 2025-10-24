package main

import (
	"fmt"

	DTO "github.org/gabrielga-dev/migratto/dto"
	migration_service "github.org/gabrielga-dev/migratto/service/migration"
)

func main() {
	fmt.Println("Welcome to Migratto!")
	migration_service.Migrate(DTO.ConfigDTO{
		DatabaseHost:     "localhost",
		DatabaseName:     "migratto",
		DatabaseUsername: "admin",
		DatabasePassword: "admin123",
		Sslmode:          "disable",
		MigrationsDir:    "./migrations",
		Log:              true,
	})
}
