package file_service

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func GetFilesFromDir(MigrationsDir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(MigrationsDir)
	if err != nil {
		return nil, fmt.Errorf("error opening migrations directory: %v", err)
	}

	sqlFiles := getOnlySqlFiles(files)

	if len(files) != len(sqlFiles) {
		return nil, fmt.Errorf("the migrations dir must have only .sql files")
	}
	return sqlFiles, nil
}

func getOnlySqlFiles(files []os.DirEntry) []os.DirEntry {
	sqlFiles := []os.DirEntry{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}
	return sqlFiles
}

func GetChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to copy file content to hasher: %w", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func GetFileTag(fileName string) string {
	return strings.Split(fileName, "_")[0]
}
