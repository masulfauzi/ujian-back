package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/config"
)

func SaveImage(file *multipart.FileHeader, folder string) (string, error) {
	cfg := config.GetUploadConfig()

	if file.Size > cfg.MaxFileSize {
		return "", fmt.Errorf("file terlalu besar, maksimal %.2f MB", float64(cfg.MaxFileSize)/1024/1024)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	ext = strings.TrimPrefix(ext, ".")

	if !isAllowedFormat(ext, cfg.AllowedFormats) {
		return "", fmt.Errorf("format file tidak diizinkan: %s", ext)
	}

	timestamp := time.Now().Unix()
	randomStr := generateRandomString(8)
	newFilename := fmt.Sprintf("%d_%s.%s", timestamp, randomStr, ext)

	uploadPath := filepath.Join(cfg.ImageUploadPath, folder)

	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat folder: %v", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %v", err)
	}
	defer src.Close()

	filePath := filepath.Join(uploadPath, newFilename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("gagal membuat file: %v", err)
	}
	defer dst.Close()

	_, err = dst.ReadFrom(src)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("gagal menyimpan file: %v", err)
	}

	return newFilename, nil
}

func DeleteImage(filename string, folder string) error {
	if filename == "" {
		return nil
	}

	cfg := config.GetUploadConfig()
	filePath := filepath.Join(cfg.ImageUploadPath, folder, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	return os.Remove(filePath)
}

func isAllowedFormat(ext string, allowed []string) bool {
	for _, a := range allowed {
		if ext == a {
			return true
		}
	}
	return false
}

func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
