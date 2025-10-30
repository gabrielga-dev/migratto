package migration_model

import (
	"os"

	file_service "github.com/gabrielga-dev/migratto/service/file"
)

type MigrationModel struct {
	Filename string
	Checksum string
	Tag      string
}

func (m MigrationModel) IsEqual(file os.DirEntry, migrationsDir string) bool {
	fileCheckSum, err := file_service.GetChecksum(migrationsDir + "/" + file.Name())
	if err != nil {
		return false
	}

	return m.Filename == file.Name() && m.Checksum == fileCheckSum
}
