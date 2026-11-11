package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eneba-intern/marketplace/backend/internal/crypto"
)

type contextKey string

const UserClaimsKey contextKey = "user_claims"

// SecurityAndCORS adds security headers and handles CORS preflight requests.
func SecurityAndCORS(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		origin := r.Header.Get("Origin")
		if origin != "" {
			allowed := false
			for _, o := range allowedOrigins {
				if o == "*" || strings.EqualFold(o, origin) {
					allowed = true
					break
				}
			}
			// Also allow localhost/127.0.0.1 or .onrender.com / .vercel.app origins for preview deployments
			if !allowed && (strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasSuffix(origin, ".onrender.com") ||
				strings.HasSuffix(origin, ".vercel.app") ||
				strings.HasSuffix(origin, ".fly.dev")) {
				allowed = true
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
				w.Header().Set("Vary", "Origin")
			}
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		start := time.Now()
		next.ServeHTTP(w, r)
		if strings.HasPrefix(r.URL.Path, "/api/") {
			log.Printf("[HTTP] %s %s (%s)", r.Method, r.URL.Path, time.Since(start).Round(time.Microsecond))
		}
	})
}

// RequireAuth validates the JWT token from either the Authorization Bearer header
// or the HTTP-only session cookie.
func RequireAuth(jwtSecret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"Authentication required. Please sign in to continue."}`))
			return
		}

		claims, err := crypto.VerifyJWT(token, jwtSecret)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"Session expired or invalid token. Please sign in again."}`))
			return
		}

		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		next(w, r.WithContext(ctx))
	}
}

// GetClaims extracts authenticated user claims from the request context.
func GetClaims(r *http.Request) *crypto.Claims {
	if val := r.Context().Value(UserClaimsKey); val != nil {
		if claims, ok := val.(*crypto.Claims); ok {
			return claims
		}
	}
	return nil
}

func extractToken(r *http.Request) string {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}
	if cookie, err := r.Cookie("eneba_session"); err == nil && cookie.Value != "" {
		return strings.TrimSpace(cookie.Value)
	}
	return ""
}
