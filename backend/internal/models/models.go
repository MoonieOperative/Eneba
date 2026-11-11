package models

import "time"

// User represents a registered marketplace user.
type User struct {
	ID            int64     `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	PasswordHash  string    `json:"-"` // Never serialized to JSON
	Role          string    `json:"role"`
	WalletBalance float64   `json:"wallet_balance"`
	CashbackEarned float64  `json:"cashback_earned"`
	AvatarColor   string    `json:"avatar_color"`
	CreatedAt     time.Time `json:"created_at"`
}

// Seller represents a verified merchant offering keys on the marketplace.
type Seller struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Verified      bool    `json:"verified"`
	RatingPercent float64 `json:"rating_percent"`
	TotalSales    int     `json:"total_sales"`
	CountryCode   string  `json:"country_code"`
	Badge         string  `json:"badge"`
	JoinedYear    int     `json:"joined_year"`
}

// Offer represents a specific seller's listing for a Product.
type Offer struct {
	ID              int64   `json:"id"`
	ProductID       int64   `json:"product_id"`
	SellerID        int64   `json:"seller_id"`
	Seller          Seller  `json:"seller"`
	Price           float64 `json:"price"`
	CashbackPercent float64 `json:"cashback_percent"`
	CashbackAmount  float64 `json:"cashback_amount"`
	StockCount      int     `json:"stock_count"`
	DeliveryTime    string  `json:"delivery_time"`
	IsFeatured      bool    `json:"is_featured"`
}

// Product represents a digital game, DLC, gift card, or subscription in the catalog.
type Product struct {
	ID               int64     `json:"id"`
	Slug             string    `json:"slug"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	Description      string    `json:"description"`
	Platform         string    `json:"platform"`          // e.g. "Steam", "Xbox Live", "PlayStation", "Nintendo eShop", "EA App", "Epic Games"
	Region           string    `json:"region"`            // e.g. "GLOBAL", "EUROPE", "UNITED STATES"
	ProductType      string    `json:"product_type"`      // e.g. "Game Key", "Gift Card", "DLC", "Subscription"
	Genre            string    `json:"genre"`             // e.g. "Action RPG", "Shooter", "Sports", "Adventure", "Strategy"
	Publisher        string    `json:"publisher"`
	Developer        string    `json:"developer"`
	ReleaseDate      string    `json:"release_date"`
	CoverImage       string    `json:"cover_image"`
	BannerImage      string    `json:"banner_image"`
	Screenshots      []string  `json:"screenshots"`
	BasePrice        float64   `json:"base_price"`        // Official MSRP / Retail price
	LowestPrice      float64   `json:"lowest_price"`      // Computed lowest active offer price
	DiscountPercent  int       `json:"discount_percent"`  // Computed discount vs BasePrice
	CashbackAmount   float64   `json:"cashback_amount"`   // Cashback on lowest offer
	Rating           float64   `json:"rating"`
	ReviewCount      int       `json:"review_count"`
	LikesCount       int       `json:"likes_count"`
	SalesCount       int       `json:"sales_count"`
	InStock          bool      `json:"in_stock"`
	OffersCount      int       `json:"offers_count"`
	Offers           []Offer   `json:"offers,omitempty"`
	SystemReqs       SystemReq `json:"system_requirements"`
	ActivationSteps  []string  `json:"activation_steps"`
	MatchScore       float64   `json:"match_score,omitempty"`
	MatchedCorrection string   `json:"matched_correction,omitempty"`
}

// SystemReq holds minimum and recommended PC/console specs.
type SystemReq struct {
	OS        string `json:"os"`
	Processor string `json:"processor"`
	Memory    string `json:"memory"`
	Graphics  string `json:"graphics"`
	Storage   string `json:"storage"`
}

// HomepageBlock represents a self-contained dynamic block rendered on the landing page.
// The frontend renders blocks strictly based on the backend response without hardcoding content.
type HomepageBlock struct {
	ID          int64                  `json:"id"`
	Slug        string                 `json:"slug"`
	BlockType   string                 `json:"block_type"` // "hero_banner", "product_carousel", "flash_deals", "category_grid", "value_props"
	Title       string                 `json:"title"`
	Subtitle    string                 `json:"subtitle"`
	AccentColor string                 `json:"accent_color"`
	BadgeText   string                 `json:"badge_text,omitempty"`
	CTAText     string                 `json:"cta_text,omitempty"`
	CTALink     string                 `json:"cta_link,omitempty"`
	SortOrder   int                    `json:"sort_order"`
	ServiceName string                 `json:"service_name"` // Identifies which backend resolver/service supplied the block
	Config      map[string]interface{} `json:"config,omitempty"`
	Products    []Product              `json:"products,omitempty"`
	Items       []BlockCustomItem      `json:"items,omitempty"`
}

// BlockCustomItem represents non-product dynamic items (e.g., category cards, trust badges, promo banners).
type BlockCustomItem struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Badge       string `json:"badge,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Image       string `json:"image,omitempty"`
	Link        string `json:"link,omitempty"`
	AccentColor string `json:"accent_color,omitempty"`
	Metric      string `json:"metric,omitempty"`
}

// CatalogFilter holds query parameters for filtering and sorting the store catalog.
type CatalogFilter struct {
	Query       string
	Platforms   []string
	Regions     []string
	Genres      []string
	ProductType string
	MinPrice    *float64
	MaxPrice    *float64
	OnlyInStock bool
	SortBy      string // "popular", "price_asc", "price_desc", "discount_desc", "rating_desc", "release_desc"
	Page        int
	PageSize    int
}

// FacetCount represents filter option counts returned dynamically by the catalog endpoint.
type FacetCount struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// CatalogResponse is returned by GET /api/v1/products.
type CatalogResponse struct {
	Products         []Product             `json:"products"`
	TotalCount       int                   `json:"total_count"`
	Page             int                   `json:"page"`
	PageSize         int                   `json:"page_size"`
	TotalPages       int                   `json:"total_pages"`
	AppliedQuery     string                `json:"applied_query,omitempty"`
	SuggestedQuery   string                `json:"suggested_query,omitempty"`
	PriceBounds      PriceBounds           `json:"price_bounds"`
	AvailableFacets  map[string][]FacetCount `json:"available_facets"`
}

// PriceBounds describes the min/max price in the catalog.
type PriceBounds struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// CheckoutItemRequest represents a line item submitted during checkout.
type CheckoutItemRequest struct {
	OfferID  int64 `json:"offer_id"`
	Quantity int   `json:"quantity"`
}

// CheckoutRequest represents the payload for POST /api/v1/checkout.
type CheckoutRequest struct {
	Items         []CheckoutItemRequest `json:"items"`
	PaymentMethod string                `json:"payment_method"` // "stripe_test_card", "paypal_sandbox", "eneba_wallet", "crypto_test"
	PromoCode     string                `json:"promo_code,omitempty"`
	CardLast4     string                `json:"card_last4,omitempty"`
}

// Order represents a completed purchase and its delivered digital keys.
type Order struct {
	ID              int64       `json:"id"`
	OrderNumber     string      `json:"order_number"`
	UserID          int64       `json:"user_id"`
	Status          string      `json:"status"`
	PaymentMethod   string      `json:"payment_method"`
	PaymentProvider string      `json:"payment_provider"`
	PaymentRef      string      `json:"payment_reference"`
	Subtotal        float64     `json:"subtotal"`
	ServiceFee      float64     `json:"service_fee"`
	DiscountAmount  float64     `json:"discount_amount"`
	TotalAmount     float64     `json:"total_amount"`
	CashbackEarned  float64     `json:"cashback_earned"`
	PromoCode       string      `json:"promo_code,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	Items           []OrderItem `json:"items"`
}

// OrderItem represents an individual purchased offer and its decrypted digital activation key.
type OrderItem struct {
	ID              int64   `json:"id"`
	OrderID         int64   `json:"order_id"`
	ProductID       int64   `json:"product_id"`
	ProductSlug     string  `json:"product_slug"`
	ProductTitle    string  `json:"product_title"`
	ProductCover    string  `json:"product_cover"`
	Platform        string  `json:"platform"`
	Region          string  `json:"region"`
	OfferID         int64   `json:"offer_id"`
	SellerName      string  `json:"seller_name"`
	UnitPrice       float64 `json:"unit_price"`
	DigitalKey      string  `json:"digital_key"` // Decrypted on demand for the authenticated buyer
	KeyFormat       string  `json:"key_format"`
	ActivationGuide string  `json:"activation_guide"`
}
