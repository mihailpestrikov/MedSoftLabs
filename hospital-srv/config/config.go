package config

import "os"

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	ServerPort      string
	MLLPPort        string
	TLSCertPath     string
	TLSKeyPath      string
	DoctorAPIURL    string
	ReceptionAPIURL string
}

func Load() *Config {
	return &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "hospital_db"),
		ServerPort:      getEnv("SERVER_PORT", "9090"),
		MLLPPort:        getEnv("MLLP_PORT", "2575"),
		TLSCertPath:     getEnv("TLS_CERT_PATH", "/app/certs/server.crt"),
		TLSKeyPath:      getEnv("TLS_KEY_PATH", "/app/certs/server.key"),
		DoctorAPIURL:    getEnv("DOCTOR_API_URL", "https://doctor-api:8081"),
		ReceptionAPIURL: getEnv("RECEPTION_API_URL", "https://reception-api:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
