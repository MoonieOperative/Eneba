export interface User {
  id: number;
  email: string;
  username: string;
  role: string;
  wallet_balance: number;
  cashback_earned: number;
  avatar_color: string;
  created_at: string;
}

export interface Seller {
  id: number;
  name: string;
  slug: string;
  verified: boolean;
  rating_percent: number;
  total_sales: number;
  country_code: string;
  badge: string;
  joined_year: number;
}

export interface Offer {
  id: number;
  product_id: number;
  seller_id: number;
  seller: Seller;
  price: number;
  cashback_percent: number;
  cashback_amount: number;
  stock_count: number;
  delivery_time: string;
  is_featured: boolean;
}

export interface SystemReq {
  os: string;
  processor: string;
  memory: string;
  graphics: string;
  storage: string;
}

export interface Product {
  id: number;
  slug: string;
  title: string;
  subtitle: string;
  description: string;
  platform: string;
  region: string;
  product_type: string;
  genre: string;
  publisher: string;
  developer: string;
  release_date: string;
  cover_image: string;
  banner_image: string;
  screenshots: string[];
  base_price: number;
  lowest_price: number;
  discount_percent: number;
  cashback_amount: number;
  rating: number;
  review_count: number;
  likes_count: number;
  sales_count: number;
  in_stock: boolean;
  offers_count: number;
  offers?: Offer[];
  system_requirements: SystemReq;
  activation_steps: string[];
  match_score?: number;
  matched_correction?: string;
}

export interface BlockCustomItem {
  title: string;
  subtitle: string;
  badge?: string;
  icon?: string;
  image?: string;
  link?: string;
  accent_color?: string;
  metric?: string;
}

export interface HomepageBlock {
  id: number;
  slug: string;
  block_type: 'hero_banner' | 'flash_deals' | 'product_grid' | 'category_grid' | 'value_props' | string;
  title: string;
  subtitle: string;
  accent_color: string;
  badge_text?: string;
  cta_text?: string;
  cta_link?: string;
  sort_order: number;
  service_name: string;
  config?: Record<string, unknown>;
  products?: Product[];
  items?: BlockCustomItem[];
}

export interface FacetCount {
  value: string;
  count: number;
}

export interface CatalogResponse {
  products: Product[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
  applied_query?: string;
  suggested_query?: string;
  price_bounds: {
    min: number;
    max: number;
  };
  available_facets: {
    platforms?: FacetCount[];
    regions?: FacetCount[];
    genres?: FacetCount[];
    product_types?: FacetCount[];
  };
}

export interface SearchSuggestResponse {
  query: string;
  suggested_query: string;
  total_matches?: number;
  results: Product[];
}

export interface CartItem {
  offerId: number;
  productId: number;
  productSlug: string;
  productTitle: string;
  productCover: string;
  platform: string;
  region: string;
  sellerName: string;
  sellerRating: number;
  unitPrice: number;
  cashbackAmount: number;
  quantity: number;
}

export interface OrderItem {
  id: number;
  order_id: number;
  product_id: number;
  product_slug: string;
  product_title: string;
  product_cover: string;
  platform: string;
  region: string;
  offer_id: number;
  seller_name: string;
  unit_price: number;
  digital_key: string;
  key_format: string;
  activation_guide: string;
}

export interface Order {
  id: number;
  order_number: string;
  user_id: number;
  status: string;
  payment_method: string;
  payment_provider: string;
  payment_reference: string;
  subtotal: number;
  service_fee: number;
  discount_amount: number;
  total_amount: number;
  cashback_earned: number;
  promo_code?: string;
  created_at: string;
  items: OrderItem[];
}
