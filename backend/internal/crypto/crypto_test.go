package crypto_test

import (
	"testing"
	"time"

	"github.com/eneba-intern/marketplace/backend/internal/crypto"
)

func TestArgon2idPasswordHashing(t *testing.T) {
	rawPass := "UltraSecret#Gaming2027!"
	hash, err := crypto.HashPassword(rawPass)
	if err != nil {
		t.Fatalf("expected no error hashing password, got: %v", err)
	}
	if hash == rawPass {
		t.Fatal("password hash must not match plaintext")
	}

	valid, err := crypto.VerifyPassword(rawPass, hash)
	if err != nil || !valid {
		t.Fatalf("expected password verification to succeed, valid=%v err=%v", valid, err)
	}

	invalid, err := crypto.VerifyPassword("WrongPassword123", hash)
	if err != nil || invalid {
		t.Fatalf("expected wrong password to fail verification")
	}
}

func TestAES256GCMKeyEncryptionAtRest(t *testing.T) {
	masterSecret := []byte("injected-32-byte-secret-key-for-aes-gcm-vault!")
	digitalKey := "STM-CYBR-94FA-8821-CC09"

	ciphertext, err := crypto.EncryptSecret(digitalKey, masterSecret)
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}
	if ciphertext == digitalKey {
		t.Fatal("ciphertext must not equal plaintext digital key")
	}

	decrypted, err := crypto.DecryptSecret(ciphertext, masterSecret)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}
	if decrypted != digitalKey {
		t.Fatalf("expected %q, got %q", digitalKey, decrypted)
	}
}

func TestJWTSignAndVerify(t *testing.T) {
	secret := []byte("test-hmac-sha256-jwt-signing-secret-key-2027")
	claims := crypto.Claims{
		UserID:   42,
		Email:    "gamer@eneba.com",
		Username: "NightCityMerc",
		Role:     "buyer",
		IssuedAt: time.Now().Unix(),
		Expires:  time.Now().Add(1 * time.Hour).Unix(),
	}

	token, err := crypto.SignJWT(claims, secret)
	if err != nil {
		t.Fatalf("SignJWT failed: %v", err)
	}

	parsed, err := crypto.VerifyJWT(token, secret)
	if err != nil {
		t.Fatalf("VerifyJWT failed: %v", err)
	}
	if parsed.UserID != 42 || parsed.Email != "gamer@eneba.com" {
		t.Fatalf("unexpected parsed claims: %+v", parsed)
	}

	if _, err := crypto.VerifyJWT(token, []byte("wrong-secret")); err == nil {
		t.Fatal("expected verification with wrong secret to fail")
	}
}
