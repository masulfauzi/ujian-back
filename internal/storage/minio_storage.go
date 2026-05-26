package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitMinioClient() error {
	cfg := config.GetMinioConfig()

	c, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("gagal inisialisasi MinIO client: %w", err)
	}

	ctx := context.Background()

	exists, err := c.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return fmt.Errorf("gagal cek bucket MinIO: %w", err)
	}

	if !exists {
		if err := c.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("gagal membuat bucket MinIO '%s': %w", cfg.Bucket, err)
		}
	}

	// Set bucket policy agar object bisa diakses publik (Direct URL)
	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[{
			"Effect":"Allow",
			"Principal":{"AWS":["*"]},
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::%s/*"]
		}]
	}`, cfg.Bucket)
	if err := c.SetBucketPolicy(ctx, cfg.Bucket, policy); err != nil {
		return fmt.Errorf("gagal set bucket policy: %w", err)
	}

	minioClient = c
	return nil
}

func UploadFile(ctx context.Context, folder string, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("gagal membuka file: %w", err)
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))

	timestamp := time.Now().Unix()
	randomStr := generateRandomString(8)
	filename := fmt.Sprintf("%d_%s%s", timestamp, randomStr, ext)
	objectName := folder + "/" + filename

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	cfg := config.GetMinioConfig()
	_, err = minioClient.PutObject(ctx, cfg.Bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("gagal upload ke MinIO: %w", err)
	}

	return filename, nil
}

func DeleteFile(ctx context.Context, folder, filename string) error {
	if filename == "" {
		return nil
	}
	cfg := config.GetMinioConfig()
	objectName := folder + "/" + filename
	return minioClient.RemoveObject(ctx, cfg.Bucket, objectName, minio.RemoveObjectOptions{})
}

func generateRandomString(length int) string {
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
