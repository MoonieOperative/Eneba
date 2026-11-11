import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import {
  CreditCard,
  Wallet,
  Lock,
  ShieldCheck,
  CheckCircle2,
  Copy,
  Check,
  Sparkles,
  KeyRound,
  Eye,
  EyeOff,
  Trash2,
  ArrowRight,
  Coins,
  UserCheck,
} from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import { marketplaceApi } from '../services/api';
import type { Order } from '../types/marketplace';

export const CheckoutPage: React.FC = () => {
  const { user, demoLogin, setUser } = useAuth();
  const { items, subtotal, totalCashback, removeFromCart, updateQuantity, clearCart, formatPrice } =
    useCart();

  const [paymentMethod, setPaymentMethod] = useState<string>('stripe_test_card');
  const [promoCode, setPromoCode] = useState<string>('INTERN2027');
  const [cardNumber, setCardNumber] = useState<string>('4242 •••• •••• 4242');
  const [cardExpiry, setCardExpiry] = useState<string>('12/29');
  const [cardCvc, setCardCvc] = useState<string>('888');

  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [completedOrder, setCompletedOrder] = useState<Order | null>(null);
  const [revealedKeys, setRevealedKeys] = useState<Record<number, boolean>>({});
  const [copiedKeyId, setCopiedKeyId] = useState<number | null>(null);

  const normalizedPromo = promoCode.trim().toUpperCase();
  const discountRate =
    normalizedPromo === 'INTERN2027' || normalizedPromo === 'ENEBA20'
      ? 0.2
      : normalizedPromo === 'ENEBA10' || normalizedPromo === 'WELCOME10'
      ? 0.1
      : 0;
  const discountAmount = Math.round(subtotal * discountRate * 100) / 100;
  const serviceFee = paymentMethod === 'eneba_wallet' ? 0 : Math.round(subtotal * 0.03 * 100) / 100;
  const finalTotal = Math.max(0, Math.round((subtotal - discountAmount + serviceFee) * 100) / 100);

  const handleCompletePayment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;
    setSubmitting(true);
    setError(null);

    try {
      const res = await marketplaceApi.checkout({
        items: items.map((i) => ({ offer_id: i.offerId, quantity: i.quantity })),
        payment_method: paymentMethod,
        promo_code: normalizedPromo,
        card_last4: cardNumber.slice(-4),
      });
      setCompletedOrder(res.order);
      if (res.user) {
        setUser(res.user);
      }
      // Reveal all keys by default on the order completion screen
      const initialReveal: Record<number, boolean> = {};
      res.order.items.forEach((item) => {
        initialReveal[item.id] = true;
      });
      setRevealedKeys(initialReveal);
      clearCart();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Checkout failed');
    } finally {
      setSubmitting(false);
    }
  };

  const handleCopyKey = (id: number, keyStr: string) => {
    void navigator.clipboard.writeText(keyStr);
    setCopiedKeyId(id);
    setTimeout(() => setCopiedKeyId(null), 2000);
  };

  // Step 3: Completed Order & Digital Key Delivery Screen
  if (completedOrder) {
    return (
      <div className="max-w-4xl mx-auto px-4 sm:px-6 py-10 space-y-8">
        <div className="rounded-3xl bg-gradient-to-br from-[#2b0e52] via-[#21093e] to-[#16052b] border-2 border-[#00e676] p-6 sm:p-10 shadow-2xl space-y-6">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-[#491d7d]">
            <div className="flex items-center gap-3.5">
              <div className="w-12 h-12 rounded-2xl bg-[#00e676]/20 border border-[#00e676] flex items-center justify-center text-[#00e676]">
                <CheckCircle2 className="w-7 h-7" />
              </div>
              <div>
                <span className="text-xs font-extrabold uppercase tracking-wider text-[#00e676]">
                  Payment Verified • Keys Decrypted from Vault
                </span>
                <h1 className="text-2xl sm:text-3xl font-black text-white">
                  Order {completedOrder.order_number}
                </h1>
              </div>
            </div>
            <div className="text-left sm:text-right text-xs text-purple-200">
              <div>Gateway: <strong className="text-white">{completedOrder.payment_provider}</strong></div>
              <div>Ref: <code className="text-[#ffd524]">{completedOrder.payment_reference}</code></div>
            </div>
          </div>

          {/* Delivered Digital Keys List */}
          <div className="space-y-4">
            <h2 className="text-sm font-extrabold uppercase tracking-wider text-[#ffd524] flex items-center gap-2">
              <KeyRound className="w-4 h-4" />
              Your Purchased Digital Activation Keys ({completedOrder.items.length})
            </h2>

            {completedOrder.items.map((item) => {
              const isRevealed = revealedKeys[item.id] ?? true;
              const maskedKey = item.digital_key.replace(/[A-Z0-9]/g, '•');

              return (
                <div
                  key={item.id}
                  className="rounded-2xl bg-[#16052b] border-2 border-[#5b259e] p-5 space-y-4"
                >
                  <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                    <div className="flex items-center gap-3.5">
                      <img
                        src={item.product_cover}
                        alt={item.product_title}
                        className="w-14 h-20 object-cover rounded-xl border border-[#562396]"
                      />
                      <div>
                        <h3 className="text-base font-extrabold text-white">{item.product_title}</h3>
                        <div className="flex items-center gap-2 mt-1">
                          <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-[#3b176b] text-purple-200">
                            {item.platform}
                          </span>
                          <span className="text-[10px] font-extrabold px-2 py-0.5 rounded bg-[#00e676]/20 text-[#00e676]">
                            {item.region}
                          </span>
                          <span className="text-xs text-purple-300">
                            Seller: <strong className="text-white">{item.seller_name}</strong>
                          </span>
                        </div>
                      </div>
                    </div>
                    <div className="text-right font-extrabold text-lg text-[#ffd524]">
                      {formatPrice(item.unit_price)}
                    </div>
                  </div>

                  {/* Decrypted Digital Key Vault Box */}
                  <div className="p-4 rounded-xl bg-[#250b45] border border-[#ffd524]/60 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                    <div>
                      <span className="text-[10px] font-extrabold uppercase tracking-wider text-purple-300 block">
                        {item.key_format}
                      </span>
                      <code className="text-lg sm:text-xl font-mono font-extrabold text-[#ffd524] tracking-wider select-all">
                        {isRevealed ? item.digital_key : maskedKey}
                      </code>
                    </div>

                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        onClick={() =>
                          setRevealedKeys((prev) => ({ ...prev, [item.id]: !isRevealed }))
                        }
                        className="px-3 py-2 rounded-lg bg-[#381366] hover:bg-[#4b1a87] text-xs font-bold text-purple-100 inline-flex items-center gap-1.5 cursor-pointer"
                      >
                        {isRevealed ? <EyeOff className="w-3.5 h-3.5" /> : <Eye className="w-3.5 h-3.5" />}
                        <span>{isRevealed ? 'Hide' : 'Reveal'}</span>
                      </button>

                      <button
                        type="button"
                        onClick={() => handleCopyKey(item.id, item.digital_key)}
                        className="px-4 py-2 rounded-lg bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] text-xs font-extrabold inline-flex items-center gap-1.5 cursor-pointer"
                      >
                        {copiedKeyId === item.id ? (
                          <>
                            <Check className="w-3.5 h-3.5" />
                            <span>Copied!</span>
                          </>
                        ) : (
                          <>
                            <Copy className="w-3.5 h-3.5" />
                            <span>Copy Key</span>
                          </>
                        )}
                      </button>
                    </div>
                  </div>

                  <p className="text-xs text-purple-300">
                    <strong className="text-purple-100">Activation Instructions:</strong> {item.activation_guide}
                  </p>
                </div>
              );
            })}
          </div>

          {/* Receipt Breakdown */}
          <div className="p-4 rounded-2xl bg-[#18062e] border border-[#3d176e] grid grid-cols-2 sm:grid-cols-4 gap-4 text-xs">
            <div>
              <span className="text-purple-400 block">Subtotal</span>
              <strong className="text-white text-sm">{formatPrice(completedOrder.subtotal)}</strong>
            </div>
            <div>
              <span className="text-purple-400 block">Promo Discount</span>
              <strong className="text-[#ff4081] text-sm">-{formatPrice(completedOrder.discount_amount)}</strong>
            </div>
            <div>
              <span className="text-purple-400 block">Total Paid</span>
              <strong className="text-[#ffd524] text-sm">{formatPrice(completedOrder.total_amount)}</strong>
            </div>
            <div>
              <span className="text-purple-400 block">Cashback Credited</span>
              <strong className="text-[#00e676] text-sm">+{formatPrice(completedOrder.cashback_earned)}</strong>
            </div>
          </div>

          <div className="flex items-center justify-between flex-wrap gap-4 pt-2">
            <Link
              to="/orders"
              className="px-6 py-3 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs uppercase tracking-wider"
            >
              View All Orders & Key Library
            </Link>
            <Link
              to="/store"
              className="px-5 py-3 rounded-xl bg-[#34125e] hover:bg-[#46197d] text-white font-bold text-xs"
            >
              Continue Shopping
            </Link>
          </div>
        </div>
      </div>
    );
  }

  // Empty Cart State
  if (items.length === 0) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <div className="rounded-3xl bg-[#240a42] border border-[#4c1e86] p-12 space-y-4">
          <h1 className="text-2xl font-black text-white">Your Shopping Cart is Empty</h1>
          <p className="text-sm text-purple-300 max-w-md mx-auto">
            Add a game key, DLC, or gift card from the catalog or any product page to proceed to checkout.
          </p>
          <Link
            to="/store"
            className="inline-flex items-center gap-2 px-6 py-3.5 rounded-xl bg-[#ffd524] text-[#16072c] font-extrabold text-xs uppercase tracking-wider"
          >
            <span>Explore Store Catalog</span>
            <ArrowRight className="w-4 h-4" />
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8">
      {/* Progress Stepper */}
      <div className="mb-8 flex items-center justify-center gap-3 text-xs font-extrabold uppercase tracking-wider">
        <span className="text-[#00e676]">1. Cart Review</span>
        <span className="text-purple-500">→</span>
        <span className="px-3 py-1 rounded-full bg-[#ffd524] text-[#16072c]">2. Checkout & Payment</span>
        <span className="text-purple-500">→</span>
        <span className="text-purple-300">3. Instant Digital Key</span>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left Column: Auth Gate (if guest) + Payment Method Selection + Cart Items */}
        <div className="lg:col-span-7 space-y-6">
          {/* Authentication Requirement Gate (matches assignment flow diagram!) */}
          {!user ? (
            <div className="rounded-3xl bg-gradient-to-r from-[#341063] to-[#240a44] border-2 border-[#ffd524] p-6 space-y-4 shadow-xl">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-xl bg-[#ffd524] text-[#16072c] flex items-center justify-center">
                  <Lock className="w-5 h-5" />
                </div>
                <div>
                  <span className="text-[11px] font-extrabold uppercase tracking-wider text-[#ffd524]">
                    Authentication Required for Checkout
                  </span>
                  <h2 className="text-lg font-black text-white">
                    Sign in to receive & store your encrypted digital keys
                  </h2>
                </div>
              </div>
              <p className="text-xs text-purple-200 leading-relaxed">
                Checkout and Order History are protected for signed-in buyers so your decrypted activation keys are permanently linked to your personal Eneba Vault.
              </p>
              <div className="flex items-center gap-3 flex-wrap">
                <button
                  type="button"
                  onClick={() => void demoLogin()}
                  className="px-5 py-3 rounded-xl bg-[#00e676] hover:bg-[#33eb91] text-[#16072c] font-black text-xs uppercase tracking-wider inline-flex items-center gap-2 cursor-pointer"
                >
                  <Sparkles className="w-4 h-4" />
                  <span>1-Click Instant Demo Buyer Sign-In (€200 Wallet)</span>
                </button>
                <Link
                  to="/login?redirect=/checkout"
                  className="px-5 py-3 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs uppercase tracking-wider"
                >
                  Sign In with Email
                </Link>
                <Link
                  to="/register?redirect=/checkout"
                  className="px-4 py-3 rounded-xl bg-[#381366] text-white font-bold text-xs"
                >
                  Create Account
                </Link>
              </div>
            </div>
          ) : (
            <div className="rounded-2xl bg-[#220940] border border-[#00e676]/40 p-4 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <UserCheck className="w-5 h-5 text-[#00e676]" />
                <div className="text-xs">
                  <span className="text-purple-300">Signed in as </span>
                  <strong className="text-white">{user.username}</strong> ({user.email})
                </div>
              </div>
              <span className="text-xs font-extrabold text-[#00e676]">
                Wallet: {formatPrice(user.wallet_balance)}
              </span>
            </div>
          )}

          {/* Payment Methods Card */}
          <div className="rounded-3xl bg-[#220940] border border-[#4a1d80] p-6 space-y-5">
            <h2 className="text-lg font-extrabold text-white flex items-center gap-2">
              <CreditCard className="w-5 h-5 text-[#ffd524]" />
              Select Payment Method (Simulated / Sandbox Mode)
            </h2>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {[
                {
                  id: 'stripe_test_card',
                  title: 'Credit / Debit Card (Stripe Test)',
                  subtitle: 'Instant Visa / Mastercard 4242 Sandbox',
                  badge: 'POPULAR',
                },
                {
                  id: 'eneba_wallet',
                  title: 'Eneba Purple Wallet',
                  subtitle: user
                    ? `Available Balance: ${formatPrice(user.wallet_balance)}`
                    : '0% Service Fee for Signed-In Users',
                  badge: 'ZERO FEE',
                },
                {
                  id: 'paypal_sandbox',
                  title: 'PayPal Express Sandbox',
                  subtitle: 'Simulated PayPal Test Checkout',
                  badge: 'SANDBOX',
                },
                {
                  id: 'crypto_test',
                  title: 'CoinGate Crypto Testnet',
                  subtitle: 'USDC / BTC / ETH Testnet Simulation',
                  badge: 'WEB3 TEST',
                },
              ].map((method) => {
                const active = paymentMethod === method.id;
                return (
                  <button
                    key={method.id}
                    type="button"
                    onClick={() => setPaymentMethod(method.id)}
                    className={`p-4 rounded-2xl border-2 text-left transition-all cursor-pointer ${
                      active
                        ? 'bg-[#31115b] border-[#ffd524] shadow-lg'
                        : 'bg-[#18062e] border-[#3d176e] hover:border-[#5e26a6]'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-xs font-extrabold text-white">{method.title}</span>
                      <span className="text-[10px] font-black px-2 py-0.5 rounded bg-[#ffd524] text-[#16072c]">
                        {method.badge}
                      </span>
                    </div>
                    <p className="text-[11px] text-purple-300">{method.subtitle}</p>
                  </button>
                );
              })}
            </div>

            {/* Simulated Card Input Fields when Stripe Test Card is selected */}
            {paymentMethod === 'stripe_test_card' && (
              <div className="p-4 rounded-2xl bg-[#17052c] border border-[#3d176e] space-y-3">
                <div className="flex items-center justify-between text-xs text-purple-300">
                  <span>Stripe Test Card (No real money charged)</span>
                  <span className="text-[#00e676] font-mono">256-bit TLS Sandbox</span>
                </div>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div className="sm:col-span-2">
                    <label className="text-[11px] text-purple-300 block mb-1">Test Card Number</label>
                    <input
                      type="text"
                      value={cardNumber}
                      onChange={(e) => setCardNumber(e.target.value)}
                      className="w-full px-3 py-2 rounded-xl bg-[#230942] border border-[#4c1e86] text-xs font-mono text-white"
                    />
                  </div>
                  <div className="grid grid-cols-2 gap-2">
                    <div>
                      <label className="text-[11px] text-purple-300 block mb-1">Expiry</label>
                      <input
                        type="text"
                        value={cardExpiry}
                        onChange={(e) => setCardExpiry(e.target.value)}
                        className="w-full px-3 py-2 rounded-xl bg-[#230942] border border-[#4c1e86] text-xs font-mono text-white"
                      />
                    </div>
                    <div>
                      <label className="text-[11px] text-purple-300 block mb-1">CVC</label>
                      <input
                        type="text"
                        value={cardCvc}
                        onChange={(e) => setCardCvc(e.target.value)}
                        className="w-full px-3 py-2 rounded-xl bg-[#230942] border border-[#4c1e86] text-xs font-mono text-white"
                      />
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Order Items Review */}
          <div className="rounded-3xl bg-[#220940] border border-[#4a1d80] p-6 space-y-4">
            <h2 className="text-base font-extrabold text-white">
              Order Items ({items.reduce((acc, i) => acc + i.quantity, 0)})
            </h2>
            <div className="divide-y divide-[#381466]">
              {items.map((item) => (
                <div key={item.offerId} className="py-3.5 flex items-center justify-between gap-4">
                  <div className="flex items-center gap-3.5 min-w-0">
                    <img
                      src={item.productCover}
                      alt={item.productTitle}
                      className="w-12 h-16 object-cover rounded-lg border border-[#562396]"
                    />
                    <div className="min-w-0">
                      <Link
                        to={`/product/${item.productSlug}`}
                        className="text-sm font-extrabold text-white hover:text-[#ffd524] truncate block"
                      >
                        {item.productTitle}
                      </Link>
                      <p className="text-xs text-purple-300">
                        {item.platform} • {item.region} • Seller: <strong className="text-white">{item.sellerName}</strong>
                      </p>
                    </div>
                  </div>

                  <div className="flex items-center gap-4 shrink-0">
                    <select
                      aria-label="Item quantity"
                      value={item.quantity}
                      onChange={(e) => updateQuantity(item.offerId, Number(e.target.value))}
                      className="bg-[#18062e] border border-[#4a1d80] rounded-lg px-2.5 py-1 text-xs font-bold text-white"
                    >
                      {[1, 2, 3, 4, 5].map((n) => (
                        <option key={n} value={n}>{n}x</option>
                      ))}
                    </select>
                    <div className="text-right">
                      <div className="text-sm font-extrabold text-[#ffd524]">
                        {formatPrice(item.unitPrice * item.quantity)}
                      </div>
                      <button
                        type="button"
                        onClick={() => removeFromCart(item.offerId)}
                        className="text-[11px] text-purple-400 hover:text-red-400 inline-flex items-center gap-1 cursor-pointer"
                      >
                        <Trash2 className="w-3 h-3" /> Remove
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Right Column: Order Summary & Pay CTA */}
        <div className="lg:col-span-5 rounded-3xl bg-[#250b45] border-2 border-[#5b259e] p-6 sm:p-8 space-y-6 sticky top-28 shadow-2xl">
          <h2 className="text-lg font-extrabold text-white">Order Summary</h2>

          {/* Promo Code Input */}
          <div className="space-y-2">
            <label className="text-xs font-bold text-purple-200 block">
              Discount / Promo Code
            </label>
            <div className="flex gap-2">
              <input
                type="text"
                value={promoCode}
                onChange={(e) => setPromoCode(e.target.value)}
                placeholder="Enter promo code (e.g. INTERN2027)"
                className="flex-1 px-3.5 py-2 rounded-xl bg-[#16062b] border border-[#4d1f85] text-xs font-mono font-bold text-[#ffd524] uppercase focus:outline-none"
              />
              <button
                type="button"
                onClick={() => setPromoCode('INTERN2027')}
                className="px-3 py-2 rounded-xl bg-[#381366] hover:bg-[#4b1a87] text-xs font-bold text-purple-100 cursor-pointer"
              >
                Use INTERN2027 (-20%)
              </button>
            </div>
          </div>

          <div className="space-y-2.5 pt-4 border-t border-[#3b156b] text-sm">
            <div className="flex justify-between text-purple-200">
              <span>Items Subtotal</span>
              <span className="font-bold text-white">{formatPrice(subtotal)}</span>
            </div>
            {discountAmount > 0 && (
              <div className="flex justify-between text-[#ff4081] font-bold">
                <span>Promo Discount ({normalizedPromo})</span>
                <span>-{formatPrice(discountAmount)}</span>
              </div>
            )}
            <div className="flex justify-between text-purple-200">
              <span>Gateway Service Fee {paymentMethod === 'eneba_wallet' && '(Waived!)'}</span>
              <span className="font-bold text-white">{formatPrice(serviceFee)}</span>
            </div>
            <div className="pt-3 border-t border-[#3b156b] flex justify-between items-baseline">
              <span className="text-base font-extrabold text-white">Total to Pay</span>
              <span className="text-3xl font-black text-[#ffd524]">{formatPrice(finalTotal)}</span>
            </div>
          </div>

          <div className="p-3 rounded-xl bg-[#00e676]/15 border border-[#00e676]/30 flex items-center justify-between text-xs font-bold text-[#00e676]">
            <span className="flex items-center gap-1.5">
              <Coins className="w-4 h-4" />
              Eneba Cashback Earned
            </span>
            <span>+{formatPrice(totalCashback)}</span>
          </div>

          {error && (
            <div className="p-3.5 rounded-xl bg-red-500/20 border border-red-500/50 text-xs text-red-200 font-semibold">
              {error}
            </div>
          )}

          <button
            type="button"
            disabled={!user || submitting}
            onClick={handleCompletePayment}
            className="w-full py-4 px-6 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] disabled:opacity-40 text-[#16072c] font-black text-sm uppercase tracking-wider flex items-center justify-center gap-2 shadow-xl shadow-[#ffd524]/25 transition-all cursor-pointer"
          >
            <Wallet className="w-4 h-4" />
            <span>
              {submitting
                ? 'Decrypting Digital Key Vault...'
                : !user
                ? 'Sign In Above to Complete Purchase'
                : `Pay ${formatPrice(finalTotal)} & Unlock Instant Key`}
            </span>
          </button>

          <div className="text-[11px] text-purple-300 text-center flex items-center justify-center gap-2">
            <ShieldCheck className="w-4 h-4 text-[#00e676]" />
            <span>Simulated Sandbox Checkout • No Real Money Involved</span>
          </div>
        </div>
      </div>
    </div>
  );
};
