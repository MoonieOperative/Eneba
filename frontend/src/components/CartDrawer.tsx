import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  X,
  Trash2,
  Plus,
  Minus,
  ShoppingBag,
  ShieldCheck,
  Zap,
  ArrowRight,
  Coins,
} from 'lucide-react';
import { useCart } from '../context/CartContext';

export const CartDrawer: React.FC = () => {
  const {
    items,
    isDrawerOpen,
    closeDrawer,
    removeFromCart,
    updateQuantity,
    subtotal,
    totalCashback,
    formatPrice,
  } = useCart();
  const navigate = useNavigate();

  if (!isDrawerOpen) return null;

  const handleProceedCheckout = () => {
    closeDrawer();
    navigate('/checkout');
  };

  return (
    <div className="fixed inset-0 z-50 flex justify-end">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/70 backdrop-blur-xs transition-opacity"
        onClick={closeDrawer}
      />

      {/* Drawer Panel */}
      <div className="relative z-10 w-full max-w-md bg-[#21093e] border-l border-[#4d1e85] h-full flex flex-col shadow-2xl">
        {/* Header */}
        <div className="px-6 py-4 bg-[#2a0c4e] border-b border-[#481d7d] flex items-center justify-between">
          <div className="flex items-center gap-2.5">
            <ShoppingBag className="w-5 h-5 text-[#ffd524]" />
            <h2 className="text-lg font-extrabold text-white">Your Shopping Cart</h2>
            <span className="px-2 py-0.5 rounded-full bg-[#ffd524] text-[#16072c] text-xs font-extrabold">
              {items.reduce((acc, i) => acc + i.quantity, 0)}
            </span>
          </div>
          <button
            type="button"
            onClick={closeDrawer}
            className="p-1.5 rounded-lg text-purple-300 hover:text-white hover:bg-[#3b156b] transition-colors cursor-pointer"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Cart Items */}
        <div className="flex-1 overflow-y-auto p-6 space-y-4">
          {items.length === 0 ? (
            <div className="text-center py-16 space-y-4">
              <div className="w-16 h-16 rounded-2xl bg-[#311159] border border-[#52208f] flex items-center justify-center mx-auto">
                <ShoppingBag className="w-8 h-8 text-purple-300" />
              </div>
              <div>
                <h3 className="text-base font-bold text-white">Your cart is empty</h3>
                <p className="text-xs text-purple-300 mt-1 max-w-xs mx-auto">
                  Explore thousands of discounted digital game keys, DLCs, and gift cards from verified merchants.
                </p>
              </div>
              <Link
                to="/store"
                onClick={closeDrawer}
                className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs uppercase tracking-wide transition-colors"
              >
                Browse Store Catalog
              </Link>
            </div>
          ) : (
            items.map((item) => (
              <div
                key={item.offerId}
                className="p-3.5 rounded-2xl bg-[#2b0f50] border border-[#4c1e86] flex gap-3.5 items-start"
              >
                <Link to={`/product/${item.productSlug}`} onClick={closeDrawer} className="shrink-0">
                  <img
                    src={item.productCover}
                    alt={item.productTitle}
                    className="w-16 h-22 object-cover rounded-xl border border-[#5d26a3]"
                  />
                </Link>

                <div className="flex-1 min-w-0">
                  <Link
                    to={`/product/${item.productSlug}`}
                    onClick={closeDrawer}
                    className="text-sm font-bold text-white hover:text-[#ffd524] line-clamp-1 transition-colors"
                  >
                    {item.productTitle}
                  </Link>

                  <div className="flex items-center gap-1.5 mt-1 flex-wrap">
                    <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-[#3d176e] text-purple-200">
                      {item.platform}
                    </span>
                    <span className="text-[10px] font-extrabold px-2 py-0.5 rounded bg-[#00e676]/20 text-[#00e676]">
                      {item.region}
                    </span>
                  </div>

                  <p className="text-[11px] text-purple-300 mt-1.5">
                    Seller: <strong className="text-purple-100">{item.sellerName}</strong> ({item.sellerRating}% ★)
                  </p>

                  <div className="mt-3 flex items-center justify-between">
                    {/* Quantity Controls */}
                    <div className="inline-flex items-center rounded-lg bg-[#1b0833] border border-[#461b7c]">
                      <button
                        type="button"
                        onClick={() => updateQuantity(item.offerId, item.quantity - 1)}
                        className="p-1.5 text-purple-200 hover:text-white cursor-pointer"
                        aria-label="Decrease quantity"
                      >
                        <Minus className="w-3.5 h-3.5" />
                      </button>
                      <span className="px-2.5 text-xs font-extrabold text-white">{item.quantity}</span>
                      <button
                        type="button"
                        onClick={() => updateQuantity(item.offerId, item.quantity + 1)}
                        className="p-1.5 text-purple-200 hover:text-white cursor-pointer"
                        aria-label="Increase quantity"
                      >
                        <Plus className="w-3.5 h-3.5" />
                      </button>
                    </div>

                    <div className="text-right">
                      <div className="text-base font-extrabold text-[#ffd524]">
                        {formatPrice(item.unitPrice * item.quantity)}
                      </div>
                      <button
                        type="button"
                        onClick={() => removeFromCart(item.offerId)}
                        className="text-[11px] text-purple-300 hover:text-red-400 inline-flex items-center gap-1 mt-0.5 cursor-pointer"
                      >
                        <Trash2 className="w-3 h-3" /> Remove
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Footer Summary & Checkout CTA */}
        {items.length > 0 && (
          <div className="p-6 bg-[#18062e] border-t border-[#481d7d] space-y-4">
            <div className="space-y-2 text-sm">
              <div className="flex items-center justify-between text-purple-200">
                <span>Cart Subtotal</span>
                <span className="font-bold text-white text-base">{formatPrice(subtotal)}</span>
              </div>
              <div className="flex items-center justify-between text-xs text-[#00e676] font-bold bg-[#00e676]/10 px-3 py-2 rounded-xl border border-[#00e676]/30">
                <span className="flex items-center gap-1.5">
                  <Coins className="w-4 h-4" />
                  Eneba Purple Cashback Reward
                </span>
                <span>+{formatPrice(totalCashback)}</span>
              </div>
            </div>

            <button
              type="button"
              onClick={handleProceedCheckout}
              className="w-full py-3.5 px-6 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-sm uppercase tracking-wider flex items-center justify-center gap-2 shadow-lg shadow-[#ffd524]/25 transition-all cursor-pointer"
            >
              <span>Proceed to Checkout</span>
              <ArrowRight className="w-4 h-4" />
            </button>

            <div className="flex items-center justify-center gap-4 text-[11px] text-purple-300">
              <span className="flex items-center gap-1">
                <Zap className="w-3.5 h-3.5 text-[#ffd524]" /> Instant Key Delivery
              </span>
              <span className="flex items-center gap-1">
                <ShieldCheck className="w-3.5 h-3.5 text-[#00e676]" /> Buyer Protection
              </span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
