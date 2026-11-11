package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eneba-intern/marketplace/backend/internal/config"
	"github.com/eneba-intern/marketplace/backend/internal/crypto"
	"github.com/eneba-intern/marketplace/backend/internal/database"
	"github.com/eneba-intern/marketplace/backend/internal/middleware"
	"github.com/eneba-intern/marketplace/backend/internal/models"
	"github.com/eneba-intern/marketplace/backend/internal/search"
)

// Server wires HTTP routes to the underlying database store and runtime configuration.
type Server struct {
	Store     *database.Store
	Config    *config.Config
	StartedAt time.Time
}

// NewServer creates a new API server instance.
func NewServer(store *database.Store, cfg *config.Config) *Server {
	// If optional DEMO_USER_EMAIL and DEMO_USER_PASSWORD were injected via environment/secrets,
	// ensure that account exists in the database.
	if cfg.DemoAccountEmail != "" && len(cfg.DemoAccountPassword) >= 8 {
		_, _ = store.CreateUser(cfg.DemoAccountEmail, "Eneba Gamer", cfg.DemoAccountPassword)
	}
	return &Server{
		Store:     store,
		Config:    cfg,
		StartedAt: time.Now().UTC(),
	}
}

// Router builds the standard library Go 1.22+ ServeMux with all API routes and SPA static fallback.
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()

	// Health & Diagnostics
	mux.HandleFunc("GET /api/v1/health", s.handleHealth)

	// Homepage Dynamic Blocks
	mux.HandleFunc("GET /api/v1/homepage/blocks", s.handleHomepageBlocks)

	// Store Catalog, Typo-Tolerant Search & Single Product Page
	mux.HandleFunc("GET /api/v1/products", s.handleListProducts)
	mux.HandleFunc("GET /api/v1/products/search/suggest", s.handleSearchSuggest)
	mux.HandleFunc("GET /api/v1/products/{slug}", s.handleGetProduct)

	// Authentication
	mux.HandleFunc("POST /api/v1/auth/register", s.handleRegister)
	mux.HandleFunc("POST /api/v1/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/v1/auth/demo", s.handleDemoAuth)
	mux.HandleFunc("POST /api/v1/auth/logout", s.handleLogout)
	mux.HandleFunc("GET /api/v1/auth/me", middleware.RequireAuth(s.Config.JWTSecret, s.handleMe))

	// Protected Checkout & Order History
	mux.HandleFunc("POST /api/v1/checkout", middleware.RequireAuth(s.Config.JWTSecret, s.handleCheckout))
	mux.HandleFunc("GET /api/v1/orders", middleware.RequireAuth(s.Config.JWTSecret, s.handleGetOrders))

	// Static SPA Frontend fallback (for unified container deployment)
	mux.HandleFunc("/", s.handleSPAFallback)

	return middleware.SecurityAndCORS(s.Config.AllowedOrigins, mux)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":         "ok",
		"service":        "eneba-marketplace-api",
		"environment":    s.Config.Environment,
		"encryption":     "AES-256-GCM + Argon2id + HMAC-SHA256",
		"uptime_seconds": int(time.Since(s.StartedAt).Seconds()),
	})
}

func (s *Server) handleHomepageBlocks(w http.ResponseWriter, r *http.Request) {
	blocks, err := s.Store.GetHomepageBlocks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to load homepage blocks")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"blocks":       blocks,
		"generated_at": time.Now().UTC(),
	})
}

func (s *Server) handleListProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := models.CatalogFilter{
		Query:       strings.TrimSpace(q.Get("q")),
		Platforms:   splitCSV(q.Get("platform")),
		Regions:     splitCSV(q.Get("region")),
		Genres:      splitCSV(q.Get("genre")),
		ProductType: strings.TrimSpace(q.Get("product_type")),
		OnlyInStock: strings.EqualFold(q.Get("in_stock"), "true") || q.Get("in_stock") == "1",
		SortBy:      strings.TrimSpace(q.Get("sort")),
		Page:        parseIntOrDefault(q.Get("page"), 1),
		PageSize:    parseIntOrDefault(q.Get("page_size"), 12),
	}

	if minStr := strings.TrimSpace(q.Get("min_price")); minStr != "" {
		if val, err := strconv.ParseFloat(minStr, 64); err == nil && val >= 0 {
			filter.MinPrice = &val
		}
	}
	if maxStr := strings.TrimSpace(q.Get("max_price")); maxStr != "" {
		if val, err := strconv.ParseFloat(maxStr, 64); err == nil && val >= 0 {
			filter.MaxPrice = &val
		}
	}

	resp, err := s.Store.QueryCatalog(filter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed querying product catalog")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleSearchSuggest(w http.ResponseWriter, r *http.Request) {
	rawQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	all, err := s.Store.LoadAllProductsWithPricing()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed loading products for search")
		return
	}

	if rawQuery == "" {
		limit := 5
		if len(all) < limit {
			limit = len(all)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"query":           "",
			"suggested_query": "",
			"results":         all[:limit],
		})
		return
	}

	ranked, suggested := search.RankProducts(all, rawQuery)
	limit := 6
	if len(ranked) < limit {
		limit = len(ranked)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"query":           rawQuery,
		"suggested_query": suggested,
		"total_matches":   len(ranked),
		"results":         ranked[:limit],
	})
}

func (s *Server) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimSpace(r.PathValue("slug"))
	if slug == "" {
		writeError(w, http.StatusBadRequest, "Product slug is required")
		return
	}

	product, err := s.Store.GetProductBySlug(slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Failed loading product details")
		return
	}

	writeJSON(w, http.StatusOK, product)
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if !strings.Contains(req.Email, "@") || len(req.Email) < 5 {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address")
		return
	}
	if len(req.Username) < 2 {
		writeError(w, http.StatusBadRequest, "Username must be at least 2 characters")
		return
	}
	if len(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	user, err := s.Store.CreateUser(req.Email, req.Username, req.Password)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}

	s.issueSessionResponse(w, http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	user, err := s.Store.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	s.issueSessionResponse(w, http.StatusOK, user)
}

// handleDemoAuth creates an ephemeral reviewer buyer account with a CSPRNG password
// at runtime so evaluators can test authenticated flows with a single click without any
// hardcoded credentials stored in the codebase.
func (s *Server) handleDemoAuth(w http.ResponseWriter, r *http.Request) {
	if s.Config.DemoAccountEmail != "" && len(s.Config.DemoAccountPassword) >= 8 {
		if user, err := s.Store.AuthenticateUser(s.Config.DemoAccountEmail, s.Config.DemoAccountPassword); err == nil {
			s.issueSessionResponse(w, http.StatusOK, user)
			return
		}
	}

	randBytes := make([]byte, 4)
	_, _ = rand.Read(randBytes)
	suffix := hex.EncodeToString(randBytes)

	passBytes := make([]byte, 16)
	_, _ = rand.Read(passBytes)
	ephemeralPassword := hex.EncodeToString(passBytes)

	email := fmt.Sprintf("player.%s@eneba.demo", suffix)
	username := fmt.Sprintf("EnebaGamer_%s", strings.ToUpper(suffix[:4]))

	user, err := s.Store.CreateUser(email, username, ephemeralPassword)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed provisioning instant demo session")
		return
	}

	s.issueSessionResponse(w, http.StatusOK, user)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "eneba_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   s.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	writeJSON(w, http.StatusOK, map[string]string{"message": "Signed out successfully"})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	user, err := s.Store.GetUserByID(claims.UserID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "User account not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
}

func (s *Server) handleCheckout(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Please sign in to complete checkout")
		return
	}

	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid checkout request payload")
		return
	}

	if req.PaymentMethod == "" {
		req.PaymentMethod = "stripe_test_card"
	}

	order, err := s.Store.ProcessCheckout(claims.UserID, req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, _ := s.Store.GetUserByID(claims.UserID)

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"order": order,
		"user":  user,
	})
}

func (s *Server) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		writeError(w, http.StatusUnauthorized, "Please sign in to view your orders")
		return
	}

	orders, err := s.Store.GetOrdersByUserID(claims.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed loading order history")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}

func (s *Server) issueSessionResponse(w http.ResponseWriter, status int, user *models.User) {
	now := time.Now().UTC()
	exp := now.Add(s.Config.TokenTTL)

	claims := crypto.Claims{
		UserID:   user.ID,
		Email:    user.Email,
		Username: user.Username,
		Role:     user.Role,
		IssuedAt: now.Unix(),
		Expires:  exp.Unix(),
	}

	token, err := crypto.SignJWT(claims, s.Config.JWTSecret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed signing authentication token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "eneba_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   s.Config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(s.Config.TokenTTL.Seconds()),
	})

	writeJSON(w, status, map[string]interface{}{
		"token":      token,
		"expires_at": exp,
		"user":       user,
	})
}

func (s *Server) handleSPAFallback(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeError(w, http.StatusNotFound, "API endpoint not found")
		return
	}

	staticDir := s.Config.StaticDir
	if staticDir == "" {
		writeJSON(w, http.StatusOK, map[string]string{
			"service": "Eneba Marketplace API",
			"docs":    "GET /api/v1/health, GET /api/v1/homepage/blocks, GET /api/v1/products",
		})
		return
	}

	cleanPath := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	if cleanPath == "." || cleanPath == "" {
		cleanPath = "index.html"
	}

	targetPath := filepath.Join(staticDir, cleanPath)
	if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
		http.ServeFile(w, r, targetPath)
		return
	}

	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); err == nil {
		http.ServeFile(w, r, indexPath)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"service": "Eneba Marketplace Backend Running",
		"health":  "/api/v1/health",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func splitCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	var out []string
	for _, part := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func parseIntOrDefault(raw string, fallback int) int {
	if val, err := strconv.Atoi(strings.TrimSpace(raw)); err == nil && val > 0 {
		return val
	}
	return fallback
}
