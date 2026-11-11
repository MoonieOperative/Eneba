# Eneba Marketplace — Intern Assignment 2027

A full-stack digital gaming marketplace web application modeled after [eneba.com](https://www.eneba.com/), featuring a **React 19 + TypeScript + Tailwind CSS v4** frontend and a **Go 1.24+ REST API + SQLite** backend.

---

## Architecture & Scope Justification

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                     React 19 + TypeScript + Tailwind v4                      │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐ ┌─────────────────────┐  │
│  │   Homepage   │ │ Store/Search │ │ Product Page │ │ Protected Checkout  │  │
│  │ (Zero Hard-  │ │ (Fuzzy Typo  │ │ (Multi-Seller│ │ & Order Key Vault   │  │
│  │  coded Data) │ │  Tolerance)  │ │  Comparison) │ │ (Argon2id + JWT)    │  │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘ └──────────┬──────────┘  │
└─────────┼────────────────┼────────────────┼────────────────────┼─────────────┘
          │                │                │                    │
          ▼                ▼                ▼                    ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                         Go REST API (/api/v1/*)                              │
│  • Dynamic Homepage Block Dispatcher (8 dedicated backend block services)    │
│  • Hybrid Fuzzy Search Engine (Damerau-Levenshtein + Trigram Jaccard)        │
│  • Secret Injection Loader (ENV + Docker/K8s *_FILE secret volume mounts)    │
│  • Crypto Engine: Argon2id Passwords, HMAC-SHA256 JWT, AES-256-GCM Key Vault │
└──────────────────────────────────────┬───────────────────────────────────────┘
                                       │
                                       ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                 SQLite Database (Pure-Go modernc.org/sqlite)                 │
│  users │ sellers │ products │ offers │ digital_keys (AES-GCM) │ orders       │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Key Engineering Decisions

1. **Zero Hardcoded Homepage Content (`GET /api/v1/homepage/blocks`)**:
   - The landing page is assembled from 8 modular blocks stored in the `homepage_blocks` table and resolved at request time by dedicated backend block services (`HeroPromotionService`, `FlashDealsService`, `BestSellersService`, `CategoryNavigationService`, `TrendingNowService`, `BudgetDealsService`, `GiftCardsAndSubsService`, `TrustAndSecurityService`).
   - The React frontend strictly renders the block structures it receives from the backend.

2. **Typo-Tolerant & Partial-Name Search (`internal/search/search.go`)**:
   - Combines token prefix matching, initialism/acronym matching (`gtav`, `bg3`), character 3-gram Jaccard similarity, and **Damerau-Levenshtein edit distance** (handling missing letters, substitutions, and adjacent transpositions).
   - Queries such as `"witcer 3"`, `"wticher"`, `"cyberpnk"`, `"eldn ring"`, `"helldver"`, or `"red ded"` resolve to the right games and return an automatic `suggested_query` banner.

3. **Multi-Seller Offer Competition (`GET /api/v1/products/{slug}`)**:
   - Every product is offered by multiple verified merchants (*KeyMaster Baltic*, *NordicGames Vault*, *PixelForge Keys*, *GameSphere EU*, *HyperLoot Digital*, *Vanguard Codes*) with distinct prices, feedback ratings, sales counts, cashback percentages, and stock counts.
   - Buyers can compare offers side-by-side and choose which seller to buy from.

4. **Security, Authentication & Secret Injection (`internal/config`, `internal/crypto`)**:
   - **No Unencrypted Credentials or Keys at Rest**:
     - Passwords are hashed using **OWASP Argon2id** (`t=2, m=64MB, p=4`) with unique 16-byte cryptographic salts and verified in constant time.
     - Digital game activation keys stored in `digital_keys` and `order_items` are encrypted at rest using **AES-256-GCM** authenticated encryption and decrypted only upon authenticated purchase/retrieval.
   - **Secret Injection (`_FILE` & Environment Variables)**:
     - Secrets (`JWT_SECRET` / `JWT_SECRET_FILE`, `KEY_ENCRYPTION_SECRET` / `KEY_ENCRYPTION_SECRET_FILE`) are injected externally via environment variables or mounted Docker/Kubernetes secret files.
     - In `APP_ENV=production`, the server fails fast at startup if secrets are not injected or are shorter than 32 bytes.
     - In `APP_ENV=development` or `test`, if no external secret is provided, ephemeral 256-bit secrets are generated in-memory via `crypto/rand` so zero static secrets ever exist in source code.

---

## Quick Start (Local Development)

### 1. Run Backend API (Go 1.24+)
```bash
cd backend
go test -v ./...
go run ./cmd/server
```
The API starts on `http://localhost:8080`.

### 2. Run Frontend Dev Server (Node 22+)
```bash
cd frontend
npm install
npm run dev
```
Open `http://localhost:5173`. Vite proxies `/api` requests to `http://localhost:8080`.

---

## Containerized Deployment (Docker & Docker Compose)

### 1. Generate Local CSPRNG Secret Files (Git-Ignored)
On Windows PowerShell:
```powershell
.\deploy\init-secrets.ps1
```
On Linux / macOS:
```bash
bash ./deploy/init-secrets.sh
```

### 2. Build & Start with Docker Compose
```bash
docker compose up --build -d
```
The unified production container serves both the React SPA and the Go API on `http://localhost:8080` with secrets mounted at `/run/secrets/jwt_secret` and `/run/secrets/key_encryption_secret`.

### 3. Deploy to Render (Free Tier)
1. Push this repository to GitHub.
2. In the [Render Dashboard](https://dashboard.render.com/), click **New +** → **Blueprint** and select your repository.
3. Render reads [`render.yaml`](./render.yaml), builds the root [`Dockerfile`](./Dockerfile), and automatically injects generated 256-bit `JWT_SECRET` and `KEY_ENCRYPTION_SECRET` values.

---

## REST API Reference

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/v1/health` | Public | Health check, uptime, and encryption status |
| `GET` | `/api/v1/homepage/blocks` | Public | Dynamic homepage blocks resolved from backend services |
| `GET` | `/api/v1/products` | Public | Store catalog with filters (`q`, `platform`, `region`, `genre`, `product_type`, `min_price`, `max_price`, `sort`, `page`) |
| `GET` | `/api/v1/products/search/suggest?q=` | Public | Instant typo-tolerant autocomplete suggestions |
| `GET` | `/api/v1/products/{slug}` | Public | Single product details + competing multi-seller offers |
| `POST` | `/api/v1/auth/register` | Public | Create buyer account (Argon2id password hashing) |
| `POST` | `/api/v1/auth/login` | Public | Authenticate buyer & issue JWT + HTTP-only cookie |
| `POST` | `/api/v1/auth/demo` | Public | Provision an ephemeral CSPRNG reviewer account for 1-click testing |
| `POST` | `/api/v1/auth/logout` | Public | Clear session cookie |
| `GET` | `/api/v1/auth/me` | Bearer / Cookie | Current buyer profile, Eneba Wallet balance, and cashback |
| `POST` | `/api/v1/checkout` | **Protected** | Process simulated payment, apply promo code (`INTERN2027`), and decrypt digital keys |
| `GET` | `/api/v1/orders` | **Protected** | Buyer order history and decrypted activation keys |
