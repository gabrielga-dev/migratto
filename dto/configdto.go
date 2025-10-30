package DTO

type ConfigDTO struct {
	DatabaseDriver   string
	DatabaseHost     string
	DatabasePort     int
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	Sslmode          string
	MigrationsDir    string
	Log              bool
}
