package config

import "os"

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	ServerPort  string
	MLLPPort    string
	TLSCertPath string
	TLSKeyPath  string
}

func Load() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "hospital_db"),
		ServerPort:  getEnv("SERVER_PORT", "8081"),
		MLLPPort:    getEnv("MLLP_PORT", "2575"),
		TLSCertPath: getEnv("TLS_CERT_PATH", "/app/certs/server.crt"),
		TLSKeyPath:  getEnv("TLS_KEY_PATH", "/app/certs/server.key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
