package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitLogger() error {

	logDir := "logs"

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return err
	}

	fileName := time.Now().Format("2006-01-02") + ".log"

	file, err := os.OpenFile(
		filepath.Join(logDir, fileName),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}

	logrus.SetOutput(file)

	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logrus.SetReportCaller(true)

	logrus.SetLevel(logrus.InfoLevel)

	gin.DefaultWriter = io.MultiWriter(file)

	return nil
}
