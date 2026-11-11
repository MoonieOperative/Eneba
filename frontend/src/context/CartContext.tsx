import React, { createContext, useContext, useEffect, useState } from 'react';
import type { CartItem, Offer, Product } from '../types/marketplace';

export type CurrencyCode = 'EUR' | 'USD' | 'GBP';

interface CartContextValue {
  items: CartItem[];
  wishlistIds: number[];
  isDrawerOpen: boolean;
  currency: CurrencyCode;
  setCurrency: (c: CurrencyCode) => void;
  formatPrice: (eurAmount: number) => string;
  openDrawer: () => void;
  closeDrawer: () => void;
  addToCart: (product: Product, offer?: Offer, openCartAfter?: boolean) => void;
  removeFromCart: (offerId: number) => void;
  updateQuantity: (offerId: number, quantity: number) => void;
  clearCart: () => void;
  toggleWishlist: (productId: number) => void;
  isWishlisted: (productId: number) => boolean;
  totalItems: number;
  subtotal: number;
  totalCashback: number;
}

const CART_STORAGE_KEY = 'eneba_cart_items_v1';
const WISHLIST_STORAGE_KEY = 'eneba_wishlist_v1';

const CURRENCY_RATES: Record<CurrencyCode, { symbol: string; rate: number }> = {
  EUR: { symbol: '€', rate: 1.0 },
  USD: { symbol: '$', rate: 1.08 },
  GBP: { symbol: '£', rate: 0.85 },
};

const CartContext = createContext<CartContextValue | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>(() => {
    try {
      const saved = localStorage.getItem(CART_STORAGE_KEY);
      return saved ? (JSON.parse(saved) as CartItem[]) : [];
    } catch {
      return [];
    }
  });

  const [wishlistIds, setWishlistIds] = useState<number[]>(() => {
    try {
      const saved = localStorage.getItem(WISHLIST_STORAGE_KEY);
      return saved ? (JSON.parse(saved) as number[]) : [];
    } catch {
      return [];
    }
  });

  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [currency, setCurrency] = useState<CurrencyCode>('EUR');

  useEffect(() => {
    try {
      localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(items));
    } catch {
      // ignore
    }
  }, [items]);

  useEffect(() => {
    try {
      localStorage.setItem(WISHLIST_STORAGE_KEY, JSON.stringify(wishlistIds));
    } catch {
      // ignore
    }
  }, [wishlistIds]);

  const formatPrice = (eurAmount: number): string => {
    const info = CURRENCY_RATES[currency];
    const converted = eurAmount * info.rate;
    return `${info.symbol}${converted.toFixed(2)}`;
  };

  const addToCart = (product: Product, offer?: Offer, openCartAfter = true) => {
    // If offer is not explicitly passed (e.g. quick-add from catalog card),
    // use the first offer or derive deterministic primary offer ID ((product.id - 1) * 4 + 1)
    const resolvedOfferId =
      offer?.id ?? (product.offers && product.offers.length > 0 ? product.offers[0].id : (product.id - 1) * 4 + 1);
    const resolvedPrice = offer?.price ?? product.lowest_price;
    const resolvedCashback = offer?.cashback_amount ?? product.cashback_amount;
    const resolvedSellerName = offer?.seller?.name ?? 'KeyMaster Baltic (Official)';
    const resolvedSellerRating = offer?.seller?.rating_percent ?? 99.6;

    setItems((prev) => {
      const existingIndex = prev.findIndex((i) => i.offerId === resolvedOfferId);
      if (existingIndex >= 0) {
        const updated = [...prev];
        updated[existingIndex] = {
          ...updated[existingIndex],
          quantity: Math.min(5, updated[existingIndex].quantity + 1),
        };
        return updated;
      }
      return [
        ...prev,
        {
          offerId: resolvedOfferId,
          productId: product.id,
          productSlug: product.slug,
          productTitle: product.title,
          productCover: product.cover_image,
          platform: product.platform,
          region: product.region,
          sellerName: resolvedSellerName,
          sellerRating: resolvedSellerRating,
          unitPrice: resolvedPrice,
          cashbackAmount: resolvedCashback,
          quantity: 1,
        },
      ];
    });

    if (openCartAfter) {
      setIsDrawerOpen(true);
    }
  };

  const removeFromCart = (offerId: number) => {
    setItems((prev) => prev.filter((i) => i.offerId !== offerId));
  };

  const updateQuantity = (offerId: number, quantity: number) => {
    if (quantity <= 0) {
      removeFromCart(offerId);
      return;
    }
    setItems((prev) =>
      prev.map((i) => (i.offerId === offerId ? { ...i, quantity: Math.min(5, quantity) } : i))
    );
  };

  const clearCart = () => setItems([]);

  const toggleWishlist = (productId: number) => {
    setWishlistIds((prev) =>
      prev.includes(productId) ? prev.filter((id) => id !== productId) : [...prev, productId]
    );
  };

  const isWishlisted = (productId: number) => wishlistIds.includes(productId);

  const totalItems = items.reduce((acc, i) => acc + i.quantity, 0);
  const subtotal = Math.round(items.reduce((acc, i) => acc + i.unitPrice * i.quantity, 0) * 100) / 100;
  const totalCashback =
    Math.round(items.reduce((acc, i) => acc + i.cashbackAmount * i.quantity, 0) * 100) / 100;

  return (
    <CartContext.Provider
      value={{
        items,
        wishlistIds,
        isDrawerOpen,
        currency,
        setCurrency,
        formatPrice,
        openDrawer: () => setIsDrawerOpen(true),
        closeDrawer: () => setIsDrawerOpen(false),
        addToCart,
        removeFromCart,
        updateQuantity,
        clearCart,
        toggleWishlist,
        isWishlisted,
        totalItems,
        subtotal,
        totalCashback,
      }}
    >
      {children}
    </CartContext.Provider>
  );
};

export function useCart(): CartContextValue {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return ctx;
}
