package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eneba-intern/marketplace/backend/internal/config"
)

func TestSecretInjectionViaFileAndProductionValidation(t *testing.T) {
	tmpDir := t.TempDir()
	jwtSecretFile := filepath.Join(tmpDir, "jwt_secret.txt")
	encSecretFile := filepath.Join(tmpDir, "enc_secret.txt")

	expectedJWT := "injected-from-mounted-secret-file-jwt-key-2027-secure"
	expectedEnc := "injected-from-mounted-secret-file-aes-key-2027-secure"

	if err := os.WriteFile(jwtSecretFile, []byte(expectedJWT+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(encSecretFile, []byte(expectedEnc+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET_FILE", jwtSecretFile)
	t.Setenv("KEY_ENCRYPTION_SECRET_FILE", encSecretFile)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected config.Load() to succeed with injected secret files, got: %v", err)
	}
	if string(cfg.JWTSecret) != expectedJWT {
		t.Fatalf("expected JWT secret %q, got %q", expectedJWT, string(cfg.JWTSecret))
	}
	if string(cfg.KeyEncryptionSecret) != expectedEnc {
		t.Fatalf("expected encryption secret %q, got %q", expectedEnc, string(cfg.KeyEncryptionSecret))
	}
}

func TestProductionFailsFastWithoutInjectedSecrets(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_SECRET_FILE", "")
	t.Setenv("KEY_ENCRYPTION_SECRET", "")
	t.Setenv("KEY_ENCRYPTION_SECRET_FILE", "")

	if _, err := config.Load(); err == nil {
		t.Fatal("expected production config load to fail fast when secrets are not injected")
	}
}
