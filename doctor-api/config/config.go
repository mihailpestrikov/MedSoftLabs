package config

import "os"

type Config struct {
	ServerPort     string
	HISHTTPAddress string
	TLSCertPath    string
	TLSKeyPath     string
}

func Load() *Config {
	return &Config{
		ServerPort:     getEnv("DOCTOR_API_PORT", "8081"),
		HISHTTPAddress: getEnv("HIS_HTTP_ADDRESS", "https://hospital-srv:9090"),
		TLSCertPath:    getEnv("TLS_CERT_PATH", "/app/certs/server.crt"),
		TLSKeyPath:     getEnv("TLS_KEY_PATH", "/app/certs/server.key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
