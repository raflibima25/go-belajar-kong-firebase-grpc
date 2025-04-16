package config

import "github.com/joho/godotenv"

type Config struct {
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	GRPCPort                string
	FirebaseCredentialsPath string
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "myapp"),
		GRPCPort:   getEnv("GRPC_PORT", "50051"),
		FirebaseCredentialsPath: getEnv("FIREBASE_CREDENTIALS", "./firebase-credentials.json"),
	}

	return config, nil
}