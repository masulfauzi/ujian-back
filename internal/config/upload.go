package config

type UploadConfig struct {
	ImageUploadPath string
	MaxFileSize     int64
	AllowedFormats  []string
	ImageBaseURL    string
}

func GetUploadConfig() UploadConfig {
	return UploadConfig{
		ImageUploadPath: "./uploads/soal",
		MaxFileSize:     5 * 1024 * 1024,
		AllowedFormats:  []string{"jpg", "jpeg", "png", "gif", "webp"},
		ImageBaseURL:    "http://localhost:3000/uploads/soal",
	}
}
