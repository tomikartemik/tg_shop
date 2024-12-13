package utils

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

func SaveFile(fileData []byte, fileName string, uploadDir string) (string, error) {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	fileExt := filepath.Ext(fileName)
	newFileName := uuid.New().String() + fileExt
	pathToFile := filepath.Join(uploadDir, newFileName)
	
	err := os.WriteFile(pathToFile, fileData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return pathToFile, nil
}
