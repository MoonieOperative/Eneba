import React, { useEffect, useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import {
  Search,
  ShoppingCart,
  Heart,
  User as UserIcon,
  LogOut,
  Package,
  Sparkles,
  Wallet,
  X,
  Globe,
  ShieldCheck,
  Zap,
  ChevronRight,
  Gamepad2,
} from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useCart, type CurrencyCode } from '../context/CartContext';
import { marketplaceApi } from '../services/api';
import type { Product } from '../types/marketplace';

const TYPO_TEST_PRESETS = [
  { label: 'Try typo: "witcer 3"', query: 'witcer 3' },
  { label: 'Try typo: "cyberpnk"', query: 'cyberpnk' },
  { label: 'Try typo: "eldn ring"', query: 'eldn ring' },
  { label: 'Try partial: "helldver"', query: 'helldver' },
  { label: 'Try acronym: "gtav"', query: 'gtav' },
];

export const Navbar: React.FC = () => {
  const { user, logout, demoLogin } = useAuth();
  const { totalItems, wishlistIds, openDrawer, currency, setCurrency, formatPrice } = useCart();
  const navigate = useNavigate();
  const location = useLocation();

  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState<Product[]>([]);
  const [suggestedCorrection, setSuggestedCorrection] = useState('');
  const [isSearchOpen, setIsSearchOpen] = useState(false);
  const [isSearching, setIsSearching] = useState(false);
  const [demoLoading, setDemoLoading] = useState(false);

  const searchContainerRef = useRef<HTMLDivElement>(null);

  // Sync search input with URL query param when on /store
  useEffect(() => {
    if (location.pathname === '/store') {
      const params = new URLSearchParams(location.search);
      const q = params.get('q') || '';
      setQuery(q);
    }
  }, [location.pathname, location.search]);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (searchContainerRef.current && !searchContainerRef.current.contains(e.target as Node)) {
        setIsSearchOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    if (!isSearchOpen) return;
    let active = true;
    setIsSearching(true);
    const timer = setTimeout(async () => {
      try {
        const res = await marketplaceApi.getSearchSuggestions(query);
        if (active) {
          setSuggestions(res.results || []);
          setSuggestedCorrection(res.suggested_query || '');
        }
      } catch {
        if (active) {
          setSuggestions([]);
        }
      } finally {
        if (active) setIsSearching(false);
      }
    }, 140);

    return () => {
      active = false;
      clearTimeout(timer);
    };
  }, [query, isSearchOpen]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSearchOpen(false);
    const trimmed = query.trim();
    if (trimmed) {
      navigate(`/store?q=${encodeURIComponent(trimmed)}`);
    } else {
      navigate('/store');
    }
  };

  const handleQuickDemoSignIn = async () => {
    setDemoLoading(true);
    try {
      await demoLogin();
    } finally {
      setDemoLoading(false);
    }
  };

  return (
    <header className="sticky top-0 z-40 w-full shadow-2xl">
      {/* Top Promotional & Trust Strip */}
      <div className="bg-[#140526] border-b border-[#3a1668] text-xs text-purple-200/90">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 h-8 flex items-center justify-between">
          <div className="flex items-center gap-4 overflow-x-auto no-scrollbar">
            <span className="inline-flex items-center gap-1.5 font-semibold text-[#ffd524]">
              <Zap className="w-3.5 h-3.5 fill-[#ffd524]" />
              INSTANT DIGITAL KEY DELIVERY
            </span>
            <span className="hidden md:inline-flex items-center gap-1.5 text-purple-300">
              <ShieldCheck className="w-3.5 h-3.5 text-[#00e676]" />
              AES-256-GCM Encrypted Key Vault
            </span>
            <span className="hidden lg:inline-flex items-center gap-1.5 text-purple-200">
              Use code <code className="px-1.5 py-0.5 rounded bg-[#3a1668] text-[#ffd524] font-bold">INTERN2027</code> at checkout for 20% OFF
            </span>
          </div>

          <div className="flex items-center gap-3 shrink-0">
            <div className="flex items-center gap-1 text-purple-200">
              <Globe className="w-3.5 h-3.5 text-purple-400" />
              <span>EN /</span>
              <select
                aria-label="Select currency"
                value={currency}
                onChange={(e) => setCurrency(e.target.value as CurrencyCode)}
                className="bg-transparent font-bold text-[#ffd524] focus:outline-none cursor-pointer"
              >
                <option value="EUR" className="bg-[#1f083d] text-white">EUR (€)</option>
                <option value="USD" className="bg-[#1f083d] text-white">USD ($)</option>
                <option value="GBP" className="bg-[#1f083d] text-white">GBP (£)</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      {/* Main Eneba Royal Purple Header */}
      <div className="bg-[#290849]/95 backdrop-blur-md border-b border-[#481d7d]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 py-3.5 flex items-center justify-between gap-4">
          {/* Eneba Brand Logo */}
          <Link to="/" className="flex items-center gap-2.5 group shrink-0">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#ffd524] via-[#ffb300] to-[#ff4081] p-0.5 shadow-lg shadow-[#ffd524]/20 group-hover:scale-105 transition-transform">
              <div className="w-full h-full bg-[#1f083d] rounded-[10px] flex items-center justify-center">
                <svg className="w-6 h-6 text-[#ffd524]" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2L14.5 9.5L22 12L14.5 14.5L12 22L9.5 14.5L2 12L9.5 9.5L12 2Z" />
                </svg>
              </div>
            </div>
            <div>
              <span className="text-2xl font-extrabold tracking-tight text-white">
                eneba
              </span>
              <span className="block text-[10px] font-bold uppercase tracking-widest text-[#ffd524] -mt-1">
                Official Store
              </span>
            </div>
          </Link>

          {/* Global Typo-Tolerant Search Bar */}
          <div ref={searchContainerRef} className="relative flex-1 max-w-2xl">
            <form onSubmit={handleSearchSubmit} className="relative flex items-center">
              <Search className="w-5 h-5 text-purple-300 absolute left-4 pointer-events-none" />
              <input
                type="text"
                value={query}
                onFocus={() => setIsSearchOpen(true)}
                onChange={(e) => {
                  setQuery(e.target.value);
                  setIsSearchOpen(true);
                }}
                placeholder="Search games, gift cards, DLCs (typo-tolerant: try 'witcer 3' or 'cyberpnk')..."
                className="w-full pl-12 pr-28 py-2.5 rounded-xl bg-[#16062b] border-2 border-[#4d1f85] focus:border-[#ffd524] text-sm text-white placeholder-purple-300/60 focus:outline-none transition-colors shadow-inner"
              />
              {query && (
                <button
                  type="button"
                  onClick={() => setQuery('')}
                  className="absolute right-20 text-purple-300 hover:text-white p-1"
                  aria-label="Clear search"
                >
                  <X className="w-4 h-4" />
                </button>
              )}
              <button
                type="submit"
                className="absolute right-1.5 px-3.5 py-1.5 rounded-lg bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs uppercase tracking-wide transition-colors cursor-pointer"
              >
                Search
              </button>
            </form>

            {/* Instant Search Suggestions & Typo-Tolerance Dropdown */}
            {isSearchOpen && (
              <div className="absolute left-0 right-0 mt-2 bg-[#220a40] border-2 border-[#5a259c] rounded-2xl shadow-2xl overflow-hidden z-50">
                {/* Quick Typo-Tolerance Test Pills */}
                <div className="px-4 py-2.5 bg-[#18062e] border-b border-[#3d176e] flex items-center gap-2 flex-wrap">
                  <span className="text-[11px] font-bold uppercase tracking-wider text-purple-300 flex items-center gap-1">
                    <Sparkles className="w-3.5 h-3.5 text-[#ffd524]" />
                    Fuzzy Search Test:
                  </span>
                  {TYPO_TEST_PRESETS.map((preset) => (
                    <button
                      key={preset.query}
                      type="button"
                      onClick={() => setQuery(preset.query)}
                      className="text-xs px-2.5 py-1 rounded-full bg-[#311159] hover:bg-[#ffd524] hover:text-[#16072c] text-purple-100 font-semibold transition-colors cursor-pointer"
                    >
                      {preset.label}
                    </button>
                  ))}
                </div>

                {/* Auto-Corrected Typo Notice */}
                {suggestedCorrection && query.trim().length > 1 && (
                  <div className="px-4 py-2 bg-[#ffd524]/15 border-b border-[#ffd524]/30 flex items-center justify-between text-xs">
                    <span className="text-purple-100">
                      Showing fuzzy matches for <strong className="text-white">"{query}"</strong> — closest match:{' '}
                      <strong className="text-[#ffd524]">{suggestedCorrection}</strong>
                    </span>
                    <span className="text-[11px] px-2 py-0.5 rounded bg-[#ffd524] text-[#16072c] font-bold">
                      TYPO TOLERANT
                    </span>
                  </div>
                )}

                {/* Results List */}
                <div className="max-h-96 overflow-y-auto divide-y divide-[#381466]">
                  {isSearching ? (
                    <div className="p-6 text-center text-sm text-purple-300">
                      Searching marketplace catalog...
                    </div>
                  ) : suggestions.length === 0 ? (
                    <div className="p-6 text-center text-sm text-purple-300">
                      No matching products found for "{query}". Try a shorter keyword or browse the full store.
                    </div>
                  ) : (
                    suggestions.map((item) => (
                      <Link
                        key={item.id}
                        to={`/product/${item.slug}`}
                        onClick={() => setIsSearchOpen(false)}
                        className="flex items-center justify-between gap-4 px-4 py-3 hover:bg-[#31125a] transition-colors group"
                      >
                        <div className="flex items-center gap-3 min-w-0">
                          <img
                            src={item.cover_image}
                            alt={item.title}
                            className="w-11 h-14 object-cover rounded-lg border border-[#52228e] shrink-0"
                          />
                          <div className="min-w-0">
                            <p className="text-sm font-bold text-white group-hover:text-[#ffd524] truncate transition-colors">
                              {item.title}
                            </p>
                            <div className="flex items-center gap-2 mt-1">
                              <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-[#3b176b] text-purple-200">
                                {item.platform}
                              </span>
                              <span className="text-[10px] font-extrabold px-2 py-0.5 rounded bg-[#00e676]/20 text-[#00e676]">
                                {item.region}
                              </span>
                              {item.matched_correction && (
                                <span className="text-[10px] font-semibold text-[#ffd524]">
                                  ✓ Fuzzy Match
                                </span>
                              )}
                            </div>
                          </div>
                        </div>

                        <div className="text-right shrink-0">
                          {item.discount_percent > 0 && (
                            <span className="inline-block text-[10px] font-extrabold px-1.5 py-0.5 rounded bg-[#ff4081] text-white mb-0.5">
                              -{item.discount_percent}%
                            </span>
                          )}
                          <div className="text-base font-extrabold text-[#ffd524]">
                            {formatPrice(item.lowest_price)}
                          </div>
                        </div>
                      </Link>
                    ))
                  )}
                </div>

                <div className="px-4 py-2.5 bg-[#18062e] border-t border-[#3d176e] flex items-center justify-between">
                  <Link
                    to={query.trim() ? `/store?q=${encodeURIComponent(query.trim())}` : '/store'}
                    onClick={() => setIsSearchOpen(false)}
                    className="text-xs font-bold text-[#ffd524] hover:underline flex items-center gap-1"
                  >
                    See all results in Store Catalog <ChevronRight className="w-4 h-4" />
                  </Link>
                  <span className="text-[11px] text-purple-400">Press Enter to search</span>
                </div>
              </div>
            )}
          </div>

          {/* Right Action Icons: Wishlist, Cart, Auth */}
          <div className="flex items-center gap-2.5 shrink-0">
            <Link
              to="/store"
              className="hidden xl:inline-flex items-center gap-1.5 px-3.5 py-2 rounded-xl bg-[#381366] hover:bg-[#491985] text-xs font-bold text-white transition-colors"
            >
              <Gamepad2 className="w-4 h-4 text-[#ffd524]" />
              Browse Store
            </Link>

            {/* Wishlist Counter */}
            <Link
              to="/store"
              title="Wishlisted Products"
              className="relative p-2.5 rounded-xl bg-[#381366] hover:bg-[#491985] text-purple-100 hover:text-white transition-colors"
            >
              <Heart className="w-5 h-5" />
              {wishlistIds.length > 0 && (
                <span className="absolute -top-1 -right-1 w-5 h-5 rounded-full bg-[#ff4081] text-white text-[10px] font-extrabold flex items-center justify-center">
                  {wishlistIds.length}
                </span>
              )}
            </Link>

            {/* Cart Button */}
            <button
              type="button"
              onClick={openDrawer}
              className="relative flex items-center gap-2 px-3.5 py-2.5 rounded-xl bg-[#381366] hover:bg-[#491985] text-white transition-colors cursor-pointer"
              aria-label="Open shopping cart"
            >
              <ShoppingCart className="w-5 h-5 text-[#ffd524]" />
              {totalItems > 0 && (
                <span className="px-2 py-0.5 rounded-full bg-[#ffd524] text-[#16072c] text-xs font-extrabold">
                  {totalItems}
                </span>
              )}
            </button>

            {/* User Profile / Auth Actions */}
            {user ? (
              <div className="flex items-center gap-2">
                <Link
                  to="/orders"
                  className="flex items-center gap-2.5 px-3 py-1.5 rounded-xl bg-[#381366] hover:bg-[#491985] border border-[#5d26a6] transition-colors"
                >
                  <div
                    className="w-7 h-7 rounded-full flex items-center justify-center text-xs font-extrabold text-[#16072c]"
                    style={{ backgroundColor: user.avatar_color || '#ffd524' }}
                  >
                    {user.username.slice(0, 2).toUpperCase()}
                  </div>
                  <div className="hidden sm:block text-left">
                    <div className="text-xs font-bold text-white leading-tight truncate max-w-[110px]">
                      {user.username}
                    </div>
                    <div className="text-[11px] font-bold text-[#00e676] flex items-center gap-1">
                      <Wallet className="w-3 h-3" />
                      {formatPrice(user.wallet_balance)}
                    </div>
                  </div>
                </Link>
                <Link
                  to="/orders"
                  title="My Orders & Digital Keys"
                  className="p-2.5 rounded-xl bg-[#381366] hover:bg-[#491985] text-[#ffd524] transition-colors"
                >
                  <Package className="w-5 h-5" />
                </Link>
                <button
                  type="button"
                  onClick={() => void logout()}
                  title="Sign Out"
                  className="p-2.5 rounded-xl bg-[#381366] hover:bg-red-500/20 text-purple-200 hover:text-red-300 transition-colors cursor-pointer"
                >
                  <LogOut className="w-5 h-5" />
                </button>
              </div>
            ) : (
              <div className="flex items-center gap-2">
                <button
                  type="button"
                  onClick={handleQuickDemoSignIn}
                  disabled={demoLoading}
                  className="hidden sm:inline-flex items-center gap-1.5 px-3 py-2 rounded-xl bg-[#00e676]/15 hover:bg-[#00e676]/25 border border-[#00e676]/40 text-[#00e676] text-xs font-extrabold transition-colors cursor-pointer"
                >
                  <Sparkles className="w-3.5 h-3.5" />
                  {demoLoading ? 'Signing in...' : '1-Click Demo Login'}
                </button>
                <Link
                  to="/login"
                  className="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs uppercase tracking-wide shadow-md shadow-[#ffd524]/10 transition-all"
                >
                  <UserIcon className="w-4 h-4" />
                  <span>Sign In</span>
                </Link>
              </div>
            )}
          </div>
        </div>

        {/* Secondary Category Navigation Bar */}
        <div className="bg-[#1f083d]/90 border-t border-[#3b1669]">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 h-10 flex items-center gap-6 overflow-x-auto text-xs font-bold text-purple-200 no-scrollbar">
            <Link to="/store" className="hover:text-[#ffd524] transition-colors shrink-0 flex items-center gap-1">
              <span>All Categories</span>
            </Link>
            <Link to="/store?platform=Steam" className="hover:text-[#ffd524] transition-colors shrink-0">
              Steam Keys
            </Link>
            <Link to="/store?platform=Xbox+Live" className="hover:text-[#ffd524] transition-colors shrink-0">
              Xbox Live & Game Pass
            </Link>
            <Link to="/store?platform=PlayStation" className="hover:text-[#ffd524] transition-colors shrink-0">
              PlayStation Network
            </Link>
            <Link to="/store?platform=Nintendo+eShop" className="hover:text-[#ffd524] transition-colors shrink-0">
              Nintendo eShop
            </Link>
            <Link to="/store?product_type=Gift+Card" className="hover:text-[#ffd524] transition-colors shrink-0">
              Gift Cards
            </Link>
            <Link to="/store?max_price=25" className="text-[#00e676] hover:text-[#ffd524] transition-colors shrink-0">
              Under €25 Deals
            </Link>
            <Link to="/store?sort=discount_desc" className="text-[#ff4081] hover:text-[#ffd524] transition-colors shrink-0">
              Flash Sales (-80%)
            </Link>
          </div>
        </div>
      </div>
    </header>
  );
};
