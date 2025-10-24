package file_service

import (
	"fmt"
	"os"
)

func GetFilesFromDir(MigrationsDir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("error opening migrations directory: %v", err)
	}
	return files, nil
}
