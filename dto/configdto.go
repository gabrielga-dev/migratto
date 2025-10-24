package DTO

type ConfigDTO struct {
	DatabaseHost     string
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	Sslmode          string
	MigrationsDir    string
	Log              bool
}
