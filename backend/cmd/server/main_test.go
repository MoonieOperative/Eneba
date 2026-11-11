package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eneba-intern/marketplace/backend/internal/config"
	"github.com/eneba-intern/marketplace/backend/internal/database"
	"github.com/eneba-intern/marketplace/backend/internal/handlers"
)

func TestFullUnifiedServerAndSPAServing(t *testing.T) {
	tmpDir := t.TempDir()
	jwtSecretPath := filepath.Join(tmpDir, "jwt.secret")
	encSecretPath := filepath.Join(tmpDir, "enc.secret")
	dbPath := filepath.Join(tmpDir, "eneba_e2e.db")

	_ = os.WriteFile(jwtSecretPath, []byte("production-injected-jwt-secret-key-2027-minimum-32-chars"), 0o600)
	_ = os.WriteFile(encSecretPath, []byte("production-injected-aes-256-gcm-secret-key-2027-32-chars"), 0o600)

	t.Setenv("APP_ENV", "production")
	t.Setenv("DB_PATH", dbPath)
	t.Setenv("STATIC_DIR", filepath.Join("..", "..", "..", "frontend", "dist"))
	t.Setenv("JWT_SECRET_FILE", jwtSecretPath)
	t.Setenv("KEY_ENCRYPTION_SECRET_FILE", encSecretPath)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load failed: %v", err)
	}

	store, err := database.New(cfg.DatabasePath, cfg.KeyEncryptionSecret)
	if err != nil {
		t.Fatalf("database.New failed: %v", err)
	}
	defer func() {
		_ = store.Close()
	}()

	srv := handlers.NewServer(store, cfg)
	router := srv.Router()

	// 1. Health endpoint
	healthReq := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	healthRec := httptest.NewRecorder()
	router.ServeHTTP(healthRec, healthReq)
	if healthRec.Code != http.StatusOK {
		t.Fatalf("health check failed: status=%d", healthRec.Code)
	}

	// 2. SPA Frontend index.html served on /store client route
	spaReq := httptest.NewRequest(http.MethodGet, "/store", nil)
	spaRec := httptest.NewRecorder()
	router.ServeHTTP(spaRec, spaReq)
	if spaRec.Code != http.StatusOK {
		t.Fatalf("SPA fallback failed: status=%d", spaRec.Code)
	}

	// 3. Instant Demo Auth + Protected Wallet Checkout
	demoReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/demo", nil)
	demoRec := httptest.NewRecorder()
	router.ServeHTTP(demoRec, demoReq)
	if demoRec.Code != http.StatusOK {
		t.Fatalf("demo auth failed: status=%d", demoRec.Code)
	}
	var demoPayload struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(demoRec.Body).Decode(&demoPayload)

	chkReq := httptest.NewRequest(http.MethodPost, "/api/v1/checkout",
		strings.NewReader(`{"items":[{"offer_id":1,"quantity":1}],"payment_method":"eneba_wallet","promo_code":"INTERN2027"}`))
	chkReq.Header.Set("Authorization", "Bearer "+demoPayload.Token)
	chkReq.Header.Set("Content-Type", "application/json")

	chkRec := httptest.NewRecorder()
	router.ServeHTTP(chkRec, chkReq)
	if chkRec.Code != http.StatusCreated {
		t.Fatalf("protected wallet checkout failed: status=%d body=%s", chkRec.Code, chkRec.Body.String())
	}
}
