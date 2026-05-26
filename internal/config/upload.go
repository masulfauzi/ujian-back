package config

import (
	"fmt"
	"os"
	"strings"
)

type UploadConfig struct {
	ImageUploadPath string
	MaxFileSize     int64
	AllowedFormats  []string
	ImageBaseURL    string
}

func GetUploadConfig() UploadConfig {
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "3000"
		}
		appURL = fmt.Sprintf("http://localhost:%s", port)
	}
	appURL = strings.TrimRight(appURL, "/")

	return UploadConfig{
		ImageUploadPath: "./uploads/soal",
		MaxFileSize:     5 * 1024 * 1024,
		AllowedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
		ImageBaseURL:    appURL + "/uploads/soal",
	}
}
