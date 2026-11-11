import type {
  CatalogResponse,
  HomepageBlock,
  Order,
  Product,
  SearchSuggestResponse,
  User,
} from '../types/marketplace';

const API_BASE = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api/v1';
const TOKEN_STORAGE_KEY = 'eneba_auth_token';

export function getStoredToken(): string | null {
  try {
    return localStorage.getItem(TOKEN_STORAGE_KEY);
  } catch {
    return null;
  }
}

export function setStoredToken(token: string | null): void {
  try {
    if (token) {
      localStorage.setItem(TOKEN_STORAGE_KEY, token);
    } else {
      localStorage.removeItem(TOKEN_STORAGE_KEY);
    }
  } catch {
    // Ignore storage errors in restricted environments
  }
}

async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getStoredToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> | undefined),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
    credentials: 'include',
  });

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    const errorMsg = (data && data.error) || `Request failed with status ${response.status}`;
    throw new Error(errorMsg);
  }
  return data as T;
}

export interface CatalogQueryParams {
  q?: string;
  platform?: string[];
  region?: string[];
  genre?: string[];
  product_type?: string;
  min_price?: string;
  max_price?: string;
  in_stock?: boolean;
  sort?: string;
  page?: number;
  page_size?: number;
}

export const marketplaceApi = {
  async getHomepageBlocks(): Promise<{ blocks: HomepageBlock[]; generated_at: string }> {
    return apiFetch('/homepage/blocks');
  },

  async getCatalog(params: CatalogQueryParams): Promise<CatalogResponse> {
    const searchParams = new URLSearchParams();
    if (params.q) searchParams.set('q', params.q);
    if (params.platform && params.platform.length > 0) {
      searchParams.set('platform', params.platform.join(','));
    }
    if (params.region && params.region.length > 0) {
      searchParams.set('region', params.region.join(','));
    }
    if (params.genre && params.genre.length > 0) {
      searchParams.set('genre', params.genre.join(','));
    }
    if (params.product_type) searchParams.set('product_type', params.product_type);
    if (params.min_price) searchParams.set('min_price', params.min_price);
    if (params.max_price) searchParams.set('max_price', params.max_price);
    if (params.in_stock) searchParams.set('in_stock', 'true');
    if (params.sort) searchParams.set('sort', params.sort);
    if (params.page) searchParams.set('page', String(params.page));
    if (params.page_size) searchParams.set('page_size', String(params.page_size));

    const qs = searchParams.toString();
    return apiFetch(`/products${qs ? `?${qs}` : ''}`);
  },

  async getSearchSuggestions(query: string): Promise<SearchSuggestResponse> {
    const qs = new URLSearchParams({ q: query }).toString();
    return apiFetch(`/products/search/suggest?${qs}`);
  },

  async getProductBySlug(slug: string): Promise<Product> {
    return apiFetch(`/products/${encodeURIComponent(slug)}`);
  },

  async register(email: string, username: string, password: string): Promise<{ token: string; user: User }> {
    const res = await apiFetch<{ token: string; user: User }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, username, password }),
    });
    setStoredToken(res.token);
    return res;
  },

  async login(email: string, password: string): Promise<{ token: string; user: User }> {
    const res = await apiFetch<{ token: string; user: User }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    setStoredToken(res.token);
    return res;
  },

  async demoLogin(): Promise<{ token: string; user: User }> {
    const res = await apiFetch<{ token: string; user: User }>('/auth/demo', {
      method: 'POST',
    });
    setStoredToken(res.token);
    return res;
  },

  async logout(): Promise<void> {
    try {
      await apiFetch('/auth/logout', { method: 'POST' });
    } finally {
      setStoredToken(null);
    }
  },

  async getMe(): Promise<{ user: User }> {
    return apiFetch('/auth/me');
  },

  async checkout(payload: {
    items: { offer_id: number; quantity: number }[];
    payment_method: string;
    promo_code?: string;
    card_last4?: string;
  }): Promise<{ order: Order; user: User }> {
    return apiFetch('/checkout', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
  },

  async getOrders(): Promise<{ orders: Order[] }> {
    return apiFetch('/orders');
  },
};
