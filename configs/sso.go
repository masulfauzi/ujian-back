package configs

// SSOConfig berisi konfigurasi login otomatis (SSO) dari aplikasi Skanida Student via sianida.
type SSOConfig struct {
	VerifyURL string
	Secret    string
}

func GetSSOConfig() *SSOConfig {
	return &SSOConfig{
		VerifyURL: getEnv("SIANIDA_SSO_VERIFY_URL", "https://apps.smkn2semarang.sch.id/api/ujian/sso-verify"),
		Secret:    getEnv("SIANIDA_SSO_SECRET", ""),
	}
}
