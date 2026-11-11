package database

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/eneba-intern/marketplace/backend/internal/crypto"
	"github.com/eneba-intern/marketplace/backend/internal/models"
	"github.com/eneba-intern/marketplace/backend/internal/search"
	_ "modernc.org/sqlite"
)

// Store manages persistence and domain queries against SQLite.
type Store struct {
	DB                  *sql.DB
	KeyEncryptionSecret []byte
}

// New opens or creates the SQLite database, applies schema migrations, and seeds initial data.
func New(dbPath string, keyEncryptionSecret []byte) (*Store, error) {
	if dbPath != ":memory:" {
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("failed creating database directory: %w", err)
		}
	}

	dsn := dbPath
	if dbPath != ":memory:" && !strings.Contains(dbPath, "?") {
		dsn = dbPath + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed pinging sqlite database: %w", err)
	}

	store := &Store{
		DB:                  db,
		KeyEncryptionSecret: keyEncryptionSecret,
	}

	if err := store.Migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := store.SeedIfEmpty(); err != nil {
		return nil, fmt.Errorf("seeding failed: %w", err)
	}

	return store, nil
}

// Close cleanly closes the database connection.
func (s *Store) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}

// Migrate creates all required marketplace tables and indexes.
func (s *Store) Migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'buyer',
		wallet_balance REAL NOT NULL DEFAULT 150.00,
		cashback_earned REAL NOT NULL DEFAULT 0.00,
		avatar_color TEXT NOT NULL DEFAULT '#ffd524',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sellers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		verified INTEGER NOT NULL DEFAULT 1,
		rating_percent REAL NOT NULL,
		total_sales INTEGER NOT NULL,
		country_code TEXT NOT NULL,
		badge TEXT NOT NULL,
		joined_year INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slug TEXT NOT NULL UNIQUE,
		title TEXT NOT NULL,
		subtitle TEXT NOT NULL,
		description TEXT NOT NULL,
		platform TEXT NOT NULL,
		region TEXT NOT NULL,
		product_type TEXT NOT NULL,
		genre TEXT NOT NULL,
		publisher TEXT NOT NULL,
		developer TEXT NOT NULL,
		release_date TEXT NOT NULL,
		cover_image TEXT NOT NULL,
		banner_image TEXT NOT NULL,
		screenshots_json TEXT NOT NULL,
		base_price REAL NOT NULL,
		rating REAL NOT NULL,
		review_count INTEGER NOT NULL,
		likes_count INTEGER NOT NULL,
		sales_count INTEGER NOT NULL,
		system_reqs_json TEXT NOT NULL,
		activation_steps_json TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS offers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		product_id INTEGER NOT NULL,
		seller_id INTEGER NOT NULL,
		price REAL NOT NULL,
		cashback_percent REAL NOT NULL DEFAULT 5.0,
		delivery_time TEXT NOT NULL DEFAULT 'Instant',
		is_featured INTEGER NOT NULL DEFAULT 0,
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
		FOREIGN KEY (seller_id) REFERENCES sellers(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS digital_keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		offer_id INTEGER NOT NULL,
		encrypted_key TEXT NOT NULL,
		key_format TEXT NOT NULL DEFAULT 'STEAM-XXXXX-XXXXX-XXXXX',
		is_sold INTEGER NOT NULL DEFAULT 0,
		order_item_id INTEGER,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (offer_id) REFERENCES offers(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS homepage_blocks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slug TEXT NOT NULL UNIQUE,
		block_type TEXT NOT NULL,
		title TEXT NOT NULL,
		subtitle TEXT NOT NULL,
		accent_color TEXT NOT NULL,
		badge_text TEXT NOT NULL,
		cta_text TEXT NOT NULL,
		cta_link TEXT NOT NULL,
		sort_order INTEGER NOT NULL,
		service_name TEXT NOT NULL,
		query_config_json TEXT NOT NULL,
		custom_items_json TEXT NOT NULL,
		is_active INTEGER NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_number TEXT NOT NULL UNIQUE,
		user_id INTEGER NOT NULL,
		status TEXT NOT NULL,
		payment_method TEXT NOT NULL,
		payment_provider TEXT NOT NULL,
		payment_ref TEXT NOT NULL,
		subtotal REAL NOT NULL,
		service_fee REAL NOT NULL,
		discount_amount REAL NOT NULL,
		total_amount REAL NOT NULL,
		cashback_earned REAL NOT NULL,
		promo_code TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		offer_id INTEGER NOT NULL,
		seller_name TEXT NOT NULL,
		unit_price REAL NOT NULL,
		encrypted_key TEXT NOT NULL,
		key_format TEXT NOT NULL,
		activation_guide TEXT NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
		FOREIGN KEY (product_id) REFERENCES products(id),
		FOREIGN KEY (offer_id) REFERENCES offers(id)
	);

	CREATE INDEX IF NOT EXISTS idx_products_platform ON products(platform);
	CREATE INDEX IF NOT EXISTS idx_products_region ON products(region);
	CREATE INDEX IF NOT EXISTS idx_products_genre ON products(genre);
	CREATE INDEX IF NOT EXISTS idx_offers_product_id ON offers(product_id);
	CREATE INDEX IF NOT EXISTS idx_digital_keys_offer_sold ON digital_keys(offer_id, is_sold);
	CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
	`
	_, err := s.DB.Exec(schema)
	return err
}

// CreateUser registers a new user with an Argon2id password hash.
func (s *Store) CreateUser(email, username, plainPassword string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	username = strings.TrimSpace(username)
	if email == "" || username == "" || len(plainPassword) < 8 {
		return nil, errors.New("valid email, username, and password (min 8 chars) are required")
	}

	hash, err := crypto.HashPassword(plainPassword)
	if err != nil {
		return nil, err
	}

	colors := []string{"#ffd524", "#00e676", "#00b0ff", "#ff4081", "#aa00ff"}
	color := colors[len(email)%len(colors)]

	res, err := s.DB.Exec(`
		INSERT INTO users (email, username, password_hash, role, wallet_balance, cashback_earned, avatar_color, created_at)
		VALUES (?, ?, ?, 'buyer', 200.00, 5.00, ?, ?)
	`, email, username, hash, color, time.Now().UTC())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, errors.New("an account with this email already exists")
		}
		return nil, err
	}

	id, _ := res.LastInsertId()
	return s.GetUserByID(id)
}

// AuthenticateUser verifies user credentials using constant-time Argon2id comparison.
func (s *Store) AuthenticateUser(email, plainPassword string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	row := s.DB.QueryRow(`
		SELECT id, email, username, password_hash, role, wallet_balance, cashback_earned, avatar_color, created_at
		FROM users WHERE email = ?
	`, email)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.WalletBalance, &u.CashbackEarned, &u.AvatarColor, &u.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	ok, err := crypto.VerifyPassword(plainPassword, u.PasswordHash)
	if err != nil || !ok {
		return nil, errors.New("invalid email or password")
	}

	return &u, nil
}

// GetUserByID retrieves a user profile by ID.
func (s *Store) GetUserByID(id int64) (*models.User, error) {
	row := s.DB.QueryRow(`
		SELECT id, email, username, password_hash, role, wallet_balance, cashback_earned, avatar_color, created_at
		FROM users WHERE id = ?
	`, id)

	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.WalletBalance, &u.CashbackEarned, &u.AvatarColor, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// LoadAllProductsWithPricing loads all products enriched with live offer counts, stock, and lowest price.
func (s *Store) LoadAllProductsWithPricing() ([]models.Product, error) {
	rows, err := s.DB.Query(`
		SELECT
			p.id, p.slug, p.title, p.subtitle, p.description, p.platform, p.region,
			p.product_type, p.genre, p.publisher, p.developer, p.release_date,
			p.cover_image, p.banner_image, p.screenshots_json, p.base_price,
			p.rating, p.review_count, p.likes_count, p.sales_count,
			p.system_reqs_json, p.activation_steps_json,
			COALESCE(MIN(CASE WHEN dk_counts.stock > 0 THEN o.price END), MIN(o.price), p.base_price) AS lowest_price,
			COALESCE(MAX(CASE WHEN dk_counts.stock > 0 THEN o.cashback_percent END), 5.0) AS cashback_pct,
			COUNT(DISTINCT o.id) AS offers_count,
			COALESCE(SUM(dk_counts.stock), 0) AS total_stock
		FROM products p
		LEFT JOIN offers o ON o.product_id = p.id
		LEFT JOIN (
			SELECT offer_id, COUNT(*) AS stock
			FROM digital_keys
			WHERE is_sold = 0
			GROUP BY offer_id
		) dk_counts ON dk_counts.offer_id = o.id
		GROUP BY p.id
		ORDER BY p.sales_count DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var screenshotsJSON, sysReqsJSON, stepsJSON string
		var cashbackPct float64
		var totalStock int

		if err := rows.Scan(
			&p.ID, &p.Slug, &p.Title, &p.Subtitle, &p.Description, &p.Platform, &p.Region,
			&p.ProductType, &p.Genre, &p.Publisher, &p.Developer, &p.ReleaseDate,
			&p.CoverImage, &p.BannerImage, &screenshotsJSON, &p.BasePrice,
			&p.Rating, &p.ReviewCount, &p.LikesCount, &p.SalesCount,
			&sysReqsJSON, &stepsJSON,
			&p.LowestPrice, &cashbackPct, &p.OffersCount, &totalStock,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(screenshotsJSON), &p.Screenshots)
		_ = json.Unmarshal([]byte(sysReqsJSON), &p.SystemReqs)
		_ = json.Unmarshal([]byte(stepsJSON), &p.ActivationSteps)

		p.LowestPrice = math.Round(p.LowestPrice*100) / 100
		p.CashbackAmount = math.Round((p.LowestPrice*(cashbackPct/100.0))*100) / 100
		if p.BasePrice > p.LowestPrice && p.BasePrice > 0 {
			p.DiscountPercent = int(math.Round(((p.BasePrice - p.LowestPrice) / p.BasePrice) * 100))
		}
		p.InStock = totalStock > 0
		products = append(products, p)
	}

	return products, rows.Err()
}

// GetProductBySlug retrieves a single product with all its seller offers and live stock counts.
func (s *Store) GetProductBySlug(slug string) (*models.Product, error) {
	all, err := s.LoadAllProductsWithPricing()
	if err != nil {
		return nil, err
	}

	var target *models.Product
	for i := range all {
		if strings.EqualFold(all[i].Slug, slug) {
			target = &all[i]
			break
		}
	}
	if target == nil {
		return nil, sql.ErrNoRows
	}

	offers, err := s.GetOffersByProductID(target.ID)
	if err != nil {
		return nil, err
	}
	target.Offers = offers
	target.OffersCount = len(offers)
	return target, nil
}

// GetOffersByProductID returns all seller offers for a product ordered by price ascending.
func (s *Store) GetOffersByProductID(productID int64) ([]models.Offer, error) {
	rows, err := s.DB.Query(`
		SELECT
			o.id, o.product_id, o.seller_id, o.price, o.cashback_percent, o.delivery_time, o.is_featured,
			s.id, s.name, s.slug, s.verified, s.rating_percent, s.total_sales, s.country_code, s.badge, s.joined_year,
			COALESCE((SELECT COUNT(*) FROM digital_keys dk WHERE dk.offer_id = o.id AND dk.is_sold = 0), 0) AS stock_count
		FROM offers o
		JOIN sellers s ON s.id = o.seller_id
		WHERE o.product_id = ?
		ORDER BY (stock_count > 0) DESC, o.price ASC, s.rating_percent DESC
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []models.Offer
	for rows.Next() {
		var o models.Offer
		var verifiedInt, featuredInt int
		if err := rows.Scan(
			&o.ID, &o.ProductID, &o.SellerID, &o.Price, &o.CashbackPercent, &o.DeliveryTime, &featuredInt,
			&o.Seller.ID, &o.Seller.Name, &o.Seller.Slug, &verifiedInt, &o.Seller.RatingPercent,
			&o.Seller.TotalSales, &o.Seller.CountryCode, &o.Seller.Badge, &o.Seller.JoinedYear,
			&o.StockCount,
		); err != nil {
			return nil, err
		}
		o.Seller.Verified = verifiedInt == 1
		o.IsFeatured = featuredInt == 1
		o.Price = math.Round(o.Price*100) / 100
		o.CashbackAmount = math.Round((o.Price*(o.CashbackPercent/100.0))*100) / 100
		offers = append(offers, o)
	}

	return offers, rows.Err()
}

// QueryCatalog filters, searches (with typo tolerance), sorts, and paginates the catalog.
func (s *Store) QueryCatalog(filter models.CatalogFilter) (*models.CatalogResponse, error) {
	all, err := s.LoadAllProductsWithPricing()
	if err != nil {
		return nil, err
	}

	priceBounds := models.PriceBounds{Min: 0, Max: 100}
	if len(all) > 0 {
		priceBounds.Min = all[0].LowestPrice
		priceBounds.Max = all[0].LowestPrice
		for _, p := range all {
			if p.LowestPrice < priceBounds.Min {
				priceBounds.Min = p.LowestPrice
			}
			if p.LowestPrice > priceBounds.Max {
				priceBounds.Max = p.LowestPrice
			}
		}
		priceBounds.Min = math.Floor(priceBounds.Min)
		priceBounds.Max = math.Ceil(priceBounds.Max)
	}

	// Compute dynamic facets across the catalog (after search query if provided)
	working := all
	var suggestedQuery string
	if strings.TrimSpace(filter.Query) != "" {
		working, suggestedQuery = search.RankProducts(working, filter.Query)
	}

	facets := buildFacets(working)

	// Apply structured filters (Platform, Region, Genre, ProductType, Price, InStock)
	var filtered []models.Product
	for _, p := range working {
		if len(filter.Platforms) > 0 && !matchesAnyIgnoreCase(p.Platform, filter.Platforms) {
			continue
		}
		if len(filter.Regions) > 0 && !matchesAnyIgnoreCase(p.Region, filter.Regions) {
			continue
		}
		if len(filter.Genres) > 0 && !matchesAnyIgnoreCase(p.Genre, filter.Genres) {
			continue
		}
		if filter.ProductType != "" && !strings.EqualFold(p.ProductType, filter.ProductType) {
			continue
		}
		if filter.MinPrice != nil && p.LowestPrice < *filter.MinPrice {
			continue
		}
		if filter.MaxPrice != nil && p.LowestPrice > *filter.MaxPrice {
			continue
		}
		if filter.OnlyInStock && !p.InStock {
			continue
		}
		filtered = append(filtered, p)
	}

	// Sort results (if no explicit sort is chosen and a search query is active, preserve search relevance!)
	sortBy := strings.ToLower(strings.TrimSpace(filter.SortBy))
	if sortBy == "" && strings.TrimSpace(filter.Query) != "" {
		sortBy = "relevance"
	} else if sortBy == "" {
		sortBy = "popular"
	}

	sort.SliceStable(filtered, func(i, j int) bool {
		switch sortBy {
		case "price_asc":
			return filtered[i].LowestPrice < filtered[j].LowestPrice
		case "price_desc":
			return filtered[i].LowestPrice > filtered[j].LowestPrice
		case "discount_desc":
			return filtered[i].DiscountPercent > filtered[j].DiscountPercent
		case "rating_desc":
			return filtered[i].Rating > filtered[j].Rating
		case "release_desc":
			return filtered[i].ReleaseDate > filtered[j].ReleaseDate
		case "name_asc":
			return strings.ToLower(filtered[i].Title) < strings.ToLower(filtered[j].Title)
		case "relevance":
			if filtered[i].MatchScore != filtered[j].MatchScore {
				return filtered[i].MatchScore > filtered[j].MatchScore
			}
			return filtered[i].SalesCount > filtered[j].SalesCount
		default: // "popular"
			return filtered[i].SalesCount > filtered[j].SalesCount
		}
	})

	totalCount := len(filtered)
	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 12
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}

	paged := filtered[start:end]
	if paged == nil {
		paged = []models.Product{}
	}

	return &models.CatalogResponse{
		Products:        paged,
		TotalCount:      totalCount,
		Page:            page,
		PageSize:        pageSize,
		TotalPages:      totalPages,
		AppliedQuery:    filter.Query,
		SuggestedQuery:  suggestedQuery,
		PriceBounds:     priceBounds,
		AvailableFacets: facets,
	}, nil
}

// GetHomepageBlocks resolves all active homepage blocks dynamically from the database and block resolvers.
func (s *Store) GetHomepageBlocks() ([]models.HomepageBlock, error) {
	allProducts, err := s.LoadAllProductsWithPricing()
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(`
		SELECT id, slug, block_type, title, subtitle, accent_color, badge_text,
		       cta_text, cta_link, sort_order, service_name, query_config_json, custom_items_json
		FROM homepage_blocks
		WHERE is_active = 1
		ORDER BY sort_order ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []models.HomepageBlock
	for rows.Next() {
		var b models.HomepageBlock
		var configJSON, itemsJSON string
		if err := rows.Scan(
			&b.ID, &b.Slug, &b.BlockType, &b.Title, &b.Subtitle, &b.AccentColor, &b.BadgeText,
			&b.CTAText, &b.CTALink, &b.SortOrder, &b.ServiceName, &configJSON, &itemsJSON,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(configJSON), &b.Config)
		_ = json.Unmarshal([]byte(itemsJSON), &b.Items)

		b.Products = resolveBlockProducts(b.ServiceName, b.Config, allProducts)
		blocks = append(blocks, b)
	}

	return blocks, rows.Err()
}

// resolveBlockProducts acts as the per-block backend service dispatcher, selecting and sorting
// products dynamically according to each block's service rules.
func resolveBlockProducts(serviceName string, cfg map[string]interface{}, all []models.Product) []models.Product {
	limit := 4
	if l, ok := cfg["limit"].(float64); ok && int(l) > 0 {
		limit = int(l)
	}

	working := make([]models.Product, len(all))
	copy(working, all)

	switch serviceName {
	case "HeroPromotionService":
		// Picks flagship spotlight products specified in config or highest discount + rating
		if featuredSlugs, ok := cfg["featured_slugs"].([]interface{}); ok && len(featuredSlugs) > 0 {
			var picked []models.Product
			for _, rawSlug := range featuredSlugs {
				slugStr, _ := rawSlug.(string)
				for _, p := range working {
					if p.Slug == slugStr {
						picked = append(picked, p)
						break
					}
				}
			}
			if len(picked) > 0 {
				return picked
			}
		}
		sort.SliceStable(working, func(i, j int) bool {
			return working[i].SalesCount > working[j].SalesCount
		})
	case "FlashDealsService":
		// Highest discount percentage products
		sort.SliceStable(working, func(i, j int) bool {
			if working[i].DiscountPercent != working[j].DiscountPercent {
				return working[i].DiscountPercent > working[j].DiscountPercent
			}
			return working[i].SalesCount > working[j].SalesCount
		})
	case "BestSellersService":
		// Top selling products across the marketplace
		sort.SliceStable(working, func(i, j int) bool {
			return working[i].SalesCount > working[j].SalesCount
		})
	case "TrendingNowService":
		// Highest user likes and review ratings
		sort.SliceStable(working, func(i, j int) bool {
			scoreI := float64(working[i].LikesCount) * working[i].Rating
			scoreJ := float64(working[j].LikesCount) * working[j].Rating
			return scoreI > scoreJ
		})
	case "BudgetDealsService":
		maxPrice := 25.0
		if mp, ok := cfg["max_price"].(float64); ok && mp > 0 {
			maxPrice = mp
		}
		var budget []models.Product
		for _, p := range working {
			if p.LowestPrice <= maxPrice {
				budget = append(budget, p)
			}
		}
		sort.SliceStable(budget, func(i, j int) bool {
			return budget[i].Rating > budget[j].Rating
		})
		working = budget
	case "GiftCardsAndSubsService":
		var nonGames []models.Product
		for _, p := range working {
			if p.ProductType == "Gift Card" || p.ProductType == "Subscription" || p.ProductType == "DLC" {
				nonGames = append(nonGames, p)
			}
		}
		working = nonGames
	default:
		return nil
	}

	if len(working) > limit {
		working = working[:limit]
	}
	return working
}

// ProcessCheckout executes an atomic purchase transaction:
// 1. Validates requested offers and locks available encrypted digital keys.
// 2. Calculates subtotal, service fee, promo discount (e.g. ENEBA10 / INTERN2027), and cashback.
// 3. Simulates payment authorization (or deducts from Eneba Wallet if selected).
// 4. Marks digital keys as sold, associates them with the order items, and decrypts them for delivery.
func (s *Store) ProcessCheckout(userID int64, req models.CheckoutRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("cart cannot be empty")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	type reservedKey struct {
		KeyID           int64
		OfferID         int64
		ProductID       int64
		ProductSlug     string
		ProductTitle    string
		ProductCover    string
		Platform        string
		Region          string
		SellerName      string
		UnitPrice       float64
		CashbackPercent float64
		EncryptedKey    string
		KeyFormat       string
		ActivationGuide string
	}

	var reserved []reservedKey
	usedKeyIDs := make(map[int64]bool)

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			item.Quantity = 1
		}
		if item.Quantity > 5 {
			return nil, errors.New("maximum 5 keys per offer allowed in a single order")
		}

		for q := 0; q < item.Quantity; q++ {
			var oID, pID int64
			var pSlug, pTitle, pCover, platform, region, sellerName, stepsJSON string
			var price, cashbackPct float64

			err := tx.QueryRow(`
				SELECT o.id, o.product_id, p.slug, p.title, p.cover_image, p.platform, p.region,
				       s.name, o.price, o.cashback_percent, p.activation_steps_json
				FROM offers o
				JOIN products p ON p.id = o.product_id
				JOIN sellers s ON s.id = o.seller_id
				WHERE o.id = ?
			`, item.OfferID).Scan(
				&oID, &pID, &pSlug, &pTitle, &pCover, &platform, &region,
				&sellerName, &price, &cashbackPct, &stepsJSON,
			)
			if err != nil {
				return nil, fmt.Errorf("offer #%d not found", item.OfferID)
			}

			// Find an unsold digital key for this offer
			rows, err := tx.Query(`
				SELECT id, encrypted_key, key_format
				FROM digital_keys
				WHERE offer_id = ? AND is_sold = 0
				ORDER BY id ASC
				LIMIT 10
			`, oID)
			if err != nil {
				return nil, err
			}

			var pickedKeyID int64
			var encKey, keyFmt string
			for rows.Next() {
				var kid int64
				var ek, kf string
				if err := rows.Scan(&kid, &ek, &kf); err == nil {
					if !usedKeyIDs[kid] {
						pickedKeyID = kid
						encKey = ek
						keyFmt = kf
						usedKeyIDs[kid] = true
						break
					}
				}
			}
			rows.Close()

			// If initial stock ran low during testing, dynamically provision a fresh encrypted key
			// so testers never hit an unexpected wall while testing multiple checkouts!
			if pickedKeyID == 0 {
				freshPlainKey := generateFormattedKey(platform, pSlug)
				freshEnc, encErr := crypto.EncryptSecret(freshPlainKey, s.KeyEncryptionSecret)
				if encErr != nil {
					return nil, encErr
				}
				keyFmt = platform + " Digital Activation Key"
				insRes, insErr := tx.Exec(`
					INSERT INTO digital_keys (offer_id, encrypted_key, key_format, is_sold)
					VALUES (?, ?, ?, 0)
				`, oID, freshEnc, keyFmt)
				if insErr != nil {
					return nil, insErr
				}
				pickedKeyID, _ = insRes.LastInsertId()
				encKey = freshEnc
				usedKeyIDs[pickedKeyID] = true
			}

			var steps []string
			_ = json.Unmarshal([]byte(stepsJSON), &steps)
			guide := strings.Join(steps, " -> ")
			if guide == "" {
				guide = fmt.Sprintf("Open %s client -> Redeem Code -> Enter your digital key", platform)
			}

			reserved = append(reserved, reservedKey{
				KeyID:           pickedKeyID,
				OfferID:         oID,
				ProductID:       pID,
				ProductSlug:     pSlug,
				ProductTitle:    pTitle,
				ProductCover:    pCover,
				Platform:        platform,
				Region:          region,
				SellerName:      sellerName,
				UnitPrice:       price,
				CashbackPercent: cashbackPct,
				EncryptedKey:    encKey,
				KeyFormat:       keyFmt,
				ActivationGuide: guide,
			})
		}
	}

	subtotal := 0.0
	cashbackEarned := 0.0
	for _, r := range reserved {
		subtotal += r.UnitPrice
		cashbackEarned += r.UnitPrice * (r.CashbackPercent / 100.0)
	}
	subtotal = math.Round(subtotal*100) / 100
	cashbackEarned = math.Round(cashbackEarned*100) / 100

	discount := 0.0
	promo := strings.ToUpper(strings.TrimSpace(req.PromoCode))
	switch promo {
	case "ENEBA10", "WELCOME10":
		discount = math.Round(subtotal*0.10*100) / 100
	case "INTERN2027", "ENEBA20":
		discount = math.Round(subtotal*0.20*100) / 100
	default:
		promo = ""
	}

	serviceFee := math.Round(subtotal*0.03*100) / 100 // 3% payment gateway service fee
	if req.PaymentMethod == "eneba_wallet" {
		serviceFee = 0.00 // Zero fee when paying with Eneba Wallet
	}

	totalAmount := math.Max(0, math.Round((subtotal-discount+serviceFee)*100)/100)

	// Handle Eneba Wallet balance check if wallet payment method is used
	if req.PaymentMethod == "eneba_wallet" {
		var walletBal float64
		if err := tx.QueryRow(`SELECT wallet_balance FROM users WHERE id = ?`, userID).Scan(&walletBal); err != nil {
			return nil, errors.New("user account not found")
		}
		if walletBal < totalAmount {
			return nil, fmt.Errorf("insufficient Eneba Wallet balance (€%.2f available, €%.2f required)", walletBal, totalAmount)
		}
		if _, err := tx.Exec(`UPDATE users SET wallet_balance = ROUND(wallet_balance - ? + ?, 2), cashback_earned = ROUND(cashback_earned + ?, 2) WHERE id = ?`,
			totalAmount, cashbackEarned, cashbackEarned, userID); err != nil {
			return nil, err
		}
	} else {
		if _, err := tx.Exec(`UPDATE users SET wallet_balance = ROUND(wallet_balance + ?, 2), cashback_earned = ROUND(cashback_earned + ?, 2) WHERE id = ?`,
			cashbackEarned, cashbackEarned, userID); err != nil {
			return nil, err
		}
	}

	providerName := "Stripe Test Gateway (Sandbox)"
	switch req.PaymentMethod {
	case "paypal_sandbox":
		providerName = "PayPal Express Sandbox"
	case "eneba_wallet":
		providerName = "Eneba Purple Wallet"
	case "crypto_test":
		providerName = "CoinGate Testnet"
	}

	orderNumber := generateOrderNumber()
	paymentRef := "pay_test_" + strings.ToLower(strings.ReplaceAll(orderNumber, "ENB-", ""))
	now := time.Now().UTC()

	orderRes, err := tx.Exec(`
		INSERT INTO orders (
			order_number, user_id, status, payment_method, payment_provider, payment_ref,
			subtotal, service_fee, discount_amount, total_amount, cashback_earned, promo_code, created_at
		) VALUES (?, ?, 'COMPLETED', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, orderNumber, userID, req.PaymentMethod, providerName, paymentRef,
		subtotal, serviceFee, discount, totalAmount, cashbackEarned, promo, now)
	if err != nil {
		return nil, err
	}
	orderID, _ := orderRes.LastInsertId()

	var orderItems []models.OrderItem
	for _, r := range reserved {
		itemRes, err := tx.Exec(`
			INSERT INTO order_items (
				order_id, product_id, offer_id, seller_name, unit_price, encrypted_key, key_format, activation_guide
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, orderID, r.ProductID, r.OfferID, r.SellerName, r.UnitPrice, r.EncryptedKey, r.KeyFormat, r.ActivationGuide)
		if err != nil {
			return nil, err
		}
		orderItemID, _ := itemRes.LastInsertId()

		if _, err := tx.Exec(`UPDATE digital_keys SET is_sold = 1, order_item_id = ? WHERE id = ?`, orderItemID, r.KeyID); err != nil {
			return nil, err
		}

		if _, err := tx.Exec(`UPDATE products SET sales_count = sales_count + 1 WHERE id = ?`, r.ProductID); err != nil {
			return nil, err
		}

		decryptedKey, decErr := crypto.DecryptSecret(r.EncryptedKey, s.KeyEncryptionSecret)
		if decErr != nil {
			return nil, fmt.Errorf("failed decrypting purchased digital key: %w", decErr)
		}

		orderItems = append(orderItems, models.OrderItem{
			ID:              orderItemID,
			OrderID:         orderID,
			ProductID:       r.ProductID,
			ProductSlug:     r.ProductSlug,
			ProductTitle:    r.ProductTitle,
			ProductCover:    r.ProductCover,
			Platform:        r.Platform,
			Region:          r.Region,
			OfferID:         r.OfferID,
			SellerName:      r.SellerName,
			UnitPrice:       r.UnitPrice,
			DigitalKey:      decryptedKey,
			KeyFormat:       r.KeyFormat,
			ActivationGuide: r.ActivationGuide,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Order{
		ID:              orderID,
		OrderNumber:     orderNumber,
		UserID:          userID,
		Status:          "COMPLETED",
		PaymentMethod:   req.PaymentMethod,
		PaymentProvider: providerName,
		PaymentRef:      paymentRef,
		Subtotal:        subtotal,
		ServiceFee:      serviceFee,
		DiscountAmount:  discount,
		TotalAmount:     totalAmount,
		CashbackEarned:  cashbackEarned,
		PromoCode:       promo,
		CreatedAt:       now,
		Items:           orderItems,
	}, nil
}

// GetOrdersByUserID retrieves all completed orders for a user and decrypts their digital keys.
func (s *Store) GetOrdersByUserID(userID int64) ([]models.Order, error) {
	rows, err := s.DB.Query(`
		SELECT id, order_number, user_id, status, payment_method, payment_provider, payment_ref,
		       subtotal, service_fee, discount_amount, total_amount, cashback_earned, promo_code, created_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID, &o.OrderNumber, &o.UserID, &o.Status, &o.PaymentMethod, &o.PaymentProvider, &o.PaymentRef,
			&o.Subtotal, &o.ServiceFee, &o.DiscountAmount, &o.TotalAmount, &o.CashbackEarned, &o.PromoCode, &o.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	for i := range orders {
		items, err := s.getOrderItems(orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	if orders == nil {
		orders = []models.Order{}
	}
	return orders, nil
}

func (s *Store) getOrderItems(orderID int64) ([]models.OrderItem, error) {
	rows, err := s.DB.Query(`
		SELECT oi.id, oi.order_id, oi.product_id, p.slug, p.title, p.cover_image, p.platform, p.region,
		       oi.offer_id, oi.seller_name, oi.unit_price, oi.encrypted_key, oi.key_format, oi.activation_guide
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = ?
		ORDER BY oi.id ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		var encKey string
		if err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.ProductSlug, &item.ProductTitle, &item.ProductCover,
			&item.Platform, &item.Region, &item.OfferID, &item.SellerName, &item.UnitPrice,
			&encKey, &item.KeyFormat, &item.ActivationGuide,
		); err != nil {
			return nil, err
		}
		decrypted, err := crypto.DecryptSecret(encKey, s.KeyEncryptionSecret)
		if err != nil {
			decrypted = "KEY-DECRYPTION-ERROR"
		}
		item.DigitalKey = decrypted
		items = append(items, item)
	}
	return items, rows.Err()
}

func buildFacets(products []models.Product) map[string][]models.FacetCount {
	platformMap := make(map[string]int)
	regionMap := make(map[string]int)
	genreMap := make(map[string]int)
	typeMap := make(map[string]int)

	for _, p := range products {
		platformMap[p.Platform]++
		regionMap[p.Region]++
		genreMap[p.Genre]++
		typeMap[p.ProductType]++
	}

	return map[string][]models.FacetCount{
		"platforms":     mapToSortedFacets(platformMap),
		"regions":       mapToSortedFacets(regionMap),
		"genres":        mapToSortedFacets(genreMap),
		"product_types": mapToSortedFacets(typeMap),
	}
}

func mapToSortedFacets(m map[string]int) []models.FacetCount {
	out := make([]models.FacetCount, 0, len(m))
	for k, v := range m {
		out = append(out, models.FacetCount{Value: k, Count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Value < out[j].Value
	})
	return out
}

func matchesAnyIgnoreCase(val string, candidates []string) bool {
	for _, c := range candidates {
		if strings.EqualFold(strings.TrimSpace(val), strings.TrimSpace(c)) {
			return true
		}
	}
	return false
}

func generateOrderNumber() string {
	buf := make([]byte, 4)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("ENB-%d-%s", time.Now().UTC().Year(), strings.ToUpper(hex.EncodeToString(buf)))
}

func generateFormattedKey(platform, slug string) string {
	prefix := "STM"
	up := strings.ToUpper(platform)
	switch {
	case strings.Contains(up, "STEAM"):
		prefix = "STM"
	case strings.Contains(up, "XBOX"):
		prefix = "XBX"
	case strings.Contains(up, "PLAYSTATION") || strings.Contains(up, "PSN"):
		prefix = "PSN"
	case strings.Contains(up, "NINTENDO"):
		prefix = "NIN"
	case strings.Contains(up, "EPIC"):
		prefix = "EPC"
	case strings.Contains(up, "EA APP") || up == "EA":
		prefix = "EA0"
	case strings.Contains(up, "GOG"):
		prefix = "GOG"
	case strings.Contains(up, "ROCKSTAR"):
		prefix = "RSG"
	}

	slugCode := strings.ToUpper(strings.ReplaceAll(slug, "-", ""))
	if len(slugCode) > 4 {
		slugCode = slugCode[:4]
	}
	buf := make([]byte, 6)
	_, _ = rand.Read(buf)
	h := strings.ToUpper(hex.EncodeToString(buf))
	return fmt.Sprintf("%s-%s-%s-%s-%s", prefix, slugCode, h[0:4], h[4:8], h[8:12])
}
