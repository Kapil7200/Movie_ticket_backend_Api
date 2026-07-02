package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Logger *log.Logger

func InitLogger() error {
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		return err
	}

	// Log file name based on current date
	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	filePath := filepath.Join("logs", fileName)

	file, err := os.OpenFile(
		filePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return err
	}

	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}