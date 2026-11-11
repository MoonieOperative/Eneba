package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eneba-intern/marketplace/backend/internal/config"
	"github.com/eneba-intern/marketplace/backend/internal/database"
	"github.com/eneba-intern/marketplace/backend/internal/handlers"
	"github.com/eneba-intern/marketplace/backend/internal/models"
)

func setupTestServer(t *testing.T) (*handlers.Server, http.Handler) {
	t.Helper()
	cfg := &config.Config{
		Environment:         "test",
		Port:                "0",
		DatabasePath:        ":memory:",
		AllowedOrigins:      []string{"*"},
		JWTSecret:           []byte("ephemeral-test-jwt-secret-key-32-bytes-minimum"),
		KeyEncryptionSecret: []byte("ephemeral-test-aes-gcm-secret-32-bytes-minimum"),
		TokenTTL:            1 * time.Hour,
	}

	store, err := database.New(":memory:", cfg.KeyEncryptionSecret)
	if err != nil {
		t.Fatalf("failed creating in-memory test database: %v", err)
	}
	t.Cleanup(func() {
		_ = store.Close()
	})

	srv := handlers.NewServer(store, cfg)
	return srv, srv.Router()
}

func TestEndToEndMarketplaceFlow(t *testing.T) {
	_, router := setupTestServer(t)

	// 1. Verify Dynamic Homepage Blocks
	t.Run("HomepageBlocksAreDynamicAndPopulated", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/homepage/blocks", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}

		var payload struct {
			Blocks []models.HomepageBlock `json:"blocks"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
			t.Fatalf("failed unmarshaling homepage blocks: %v", err)
		}
		if len(payload.Blocks) < 6 {
			t.Fatalf("expected at least 6 dynamic homepage blocks, got %d", len(payload.Blocks))
		}
		if len(payload.Blocks[0].Products) == 0 {
			t.Fatalf("expected hero block to contain resolved products")
		}
	})

	// 2. Verify Store Catalog Filtering & Typo-Tolerant Search
	t.Run("CatalogSearchAndFilters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products?q=cyberpnk&platform=Steam", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}

		var resp models.CatalogResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed decoding catalog response: %v", err)
		}
		if len(resp.Products) == 0 {
			t.Fatal("expected typo-tolerant search for 'cyberpnk' to return Cyberpunk 2077")
		}
		if resp.Products[0].Slug != "cyberpunk-2077-ultimate-edition-steam-key-global" {
			t.Fatalf("unexpected top search match: %s", resp.Products[0].Slug)
		}
	})

	// 3. Verify Single Product Page & Multi-Seller Offers
	var selectedOfferID int64
	t.Run("SingleProductPageWithMultiSellerOffers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/cyberpunk-2077-ultimate-edition-steam-key-global", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rec.Code)
		}

		var prod models.Product
		if err := json.Unmarshal(rec.Body.Bytes(), &prod); err != nil {
			t.Fatalf("failed decoding product: %v", err)
		}
		if len(prod.Offers) < 3 {
			t.Fatalf("expected multiple competing seller offers, got %d", len(prod.Offers))
		}
		selectedOfferID = prod.Offers[0].ID
	})

	// 4. Verify Protected Checkout Rejects Unauthenticated Requests
	t.Run("UnauthenticatedCheckoutRejected", func(t *testing.T) {
		body := `{"items":[{"offer_id":1,"quantity":1}],"payment_method":"stripe_test_card"}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 Unauthorized for guest checkout, got %d", rec.Code)
		}
	})

	// 5. Register User, Authenticate, Checkout & Receive Decrypted Digital Key
	t.Run("RegisterAndCompleteCheckoutWithKeyDelivery", func(t *testing.T) {
		regBody := `{"email":"tester@eneba.com","username":"VaultHunter","password":"StrongPassword!2027"}`
		regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
		regRec := httptest.NewRecorder()
		router.ServeHTTP(regRec, regReq)

		if regRec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created on registration, got %d: %s", regRec.Code, regRec.Body.String())
		}

		var authResp struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}
		if err := json.Unmarshal(regRec.Body.Bytes(), &authResp); err != nil {
			t.Fatalf("failed decoding auth response: %v", err)
		}
		if authResp.Token == "" {
			t.Fatal("expected JWT token in registration response")
		}

		checkoutPayload, _ := json.Marshal(models.CheckoutRequest{
			Items:         []models.CheckoutItemRequest{{OfferID: selectedOfferID, Quantity: 1}},
			PaymentMethod: "stripe_test_card",
			PromoCode:     "INTERN2027",
		})
		chkReq := httptest.NewRequest(http.MethodPost, "/api/v1/checkout", bytes.NewReader(checkoutPayload))
		chkReq.Header.Set("Authorization", "Bearer "+authResp.Token)
		chkRec := httptest.NewRecorder()
		router.ServeHTTP(chkRec, chkReq)

		if chkRec.Code != http.StatusCreated {
			t.Fatalf("expected 201 Created on checkout, got %d: %s", chkRec.Code, chkRec.Body.String())
		}

		var chkResp struct {
			Order models.Order `json:"order"`
		}
		if err := json.Unmarshal(chkRec.Body.Bytes(), &chkResp); err != nil {
			t.Fatalf("failed decoding checkout response: %v", err)
		}

		if len(chkResp.Order.Items) != 1 {
			t.Fatalf("expected 1 order item, got %d", len(chkResp.Order.Items))
		}
		deliveredKey := chkResp.Order.Items[0].DigitalKey
		if len(deliveredKey) < 10 {
			t.Fatalf("expected decrypted digital key, got %q", deliveredKey)
		}
		if chkResp.Order.DiscountAmount <= 0 {
			t.Fatalf("expected INTERN2027 promo code discount to be applied")
		}

		// Verify Order History Endpoint
		ordReq := httptest.NewRequest(http.MethodGet, "/api/v1/orders", nil)
		ordReq.Header.Set("Authorization", "Bearer "+authResp.Token)
		ordRec := httptest.NewRecorder()
		router.ServeHTTP(ordRec, ordReq)

		if ordRec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK on /api/v1/orders, got %d", ordRec.Code)
		}
	})
}
