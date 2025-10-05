package config

import (
	"os"
)

type Config struct {
	ServerPort       string
	ACKListenerPort  string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	HISAddress       string
	TLSCertPath      string
	TLSKeyPath       string
}

func Load() *Config {
	return &Config{
		ServerPort:       getEnv("RECEPTION_API_PORT", "8080"),
		ACKListenerPort:  getEnv("ACK_LISTENER_PORT", "2576"),
		DBHost:           getEnv("RECEPTION_DB_HOST", "localhost"),
		DBPort:           getEnv("RECEPTION_DB_PORT", "5432"),
		DBUser:           getEnv("RECEPTION_DB_USER", "reception_user"),
		DBPassword:       getEnv("RECEPTION_DB_PASSWORD", "reception_password"),
		DBName:           getEnv("RECEPTION_DB_NAME", "reception_db"),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret-key-to-change"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_TOKEN_DURATION", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_TOKEN_DURATION", "168h"),
		HISAddress:       getEnv("HIS_MLLP_ADDRESS", "localhost:2575"),
		TLSCertPath:      getEnv("TLS_CERT_PATH", "../certs/server.crt"),
		TLSKeyPath:       getEnv("TLS_KEY_PATH", "../certs/server.key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
