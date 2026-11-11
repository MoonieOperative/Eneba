package config

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Config holds the runtime configuration loaded exclusively from environment
// variables or mounted secret files (Docker/Kubernetes _FILE convention).
// No static credentials or secrets are stored in source code.
type Config struct {
	Environment         string
	Port                string
	DatabasePath        string
	StaticDir           string
	AllowedOrigins      []string
	JWTSecret           []byte
	KeyEncryptionSecret []byte
	TokenTTL            time.Duration
	CookieSecure        bool
	DemoAccountEmail    string
	DemoAccountPassword string
}

// Load reads configuration and injected secrets from the environment.
func Load() (*Config, error) {
	env := strings.ToLower(strings.TrimSpace(getEnvOrDefault("APP_ENV", "development")))
	port := getEnvOrDefault("PORT", "8080")
	dbPath := getEnvOrDefault("DB_PATH", "./data/eneba.db")
	staticDir := getEnvOrDefault("STATIC_DIR", "../frontend/dist")

	originsRaw := getEnvOrDefault("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000,http://127.0.0.1:5173,http://localhost:8080")
	var origins []string
	for _, o := range strings.Split(originsRaw, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	jwtSecretStr, err := readSecret("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	encSecretStr, err := readSecret("KEY_ENCRYPTION_SECRET")
	if err != nil {
		return nil, err
	}

	if env == "production" {
		if len(jwtSecretStr) < 32 {
			return nil, errors.New("security error: JWT_SECRET (or JWT_SECRET_FILE) must be injected and at least 32 characters in production")
		}
		if len(encSecretStr) < 32 {
			return nil, errors.New("security error: KEY_ENCRYPTION_SECRET (or KEY_ENCRYPTION_SECRET_FILE) must be injected and at least 32 characters in production")
		}
	} else {
		if len(jwtSecretStr) < 32 {
			ephemeral, genErr := generateRandomHex(32)
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate ephemeral JWT secret: %w", genErr)
			}
			jwtSecretStr = ephemeral
			log.Printf("[SECURITY] No JWT_SECRET injected; generated ephemeral CSPRNG secret for %s session", env)
		}
		if len(encSecretStr) < 32 {
			ephemeral, genErr := generateRandomHex(32)
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate ephemeral encryption secret: %w", genErr)
			}
			encSecretStr = ephemeral
			log.Printf("[SECURITY] No KEY_ENCRYPTION_SECRET injected; generated ephemeral CSPRNG AES-256 secret for %s session", env)
		}
	}

	demoEmail, _ := readSecret("DEMO_USER_EMAIL")
	demoPassword, _ := readSecret("DEMO_USER_PASSWORD")

	ttlHours := 24 * time.Hour
	cookieSecure := env == "production" || strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")

	return &Config{
		Environment:         env,
		Port:                port,
		DatabasePath:        dbPath,
		StaticDir:           staticDir,
		AllowedOrigins:      origins,
		JWTSecret:           []byte(jwtSecretStr),
		KeyEncryptionSecret: []byte(encSecretStr),
		TokenTTL:            ttlHours,
		CookieSecure:        cookieSecure,
		DemoAccountEmail:    demoEmail,
		DemoAccountPassword: demoPassword,
	}, nil
}

// readSecret resolves a secret from <NAME>_FILE (container secret mount) first,
// then falls back to the <NAME> environment variable.
func readSecret(name string) (string, error) {
	fileEnv := name + "_FILE"
	if filePath := strings.TrimSpace(os.Getenv(fileEnv)); filePath != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed reading secret file %s from %s: %w", fileEnv, filePath, err)
		}
		return strings.TrimSpace(string(data)), nil
	}
	return strings.TrimSpace(os.Getenv(name)), nil
}

func getEnvOrDefault(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func generateRandomHex(numBytes int) (string, error) {
	buf := make([]byte, numBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
