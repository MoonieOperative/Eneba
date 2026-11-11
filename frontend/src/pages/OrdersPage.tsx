import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  Package,
  KeyRound,
  Copy,
  Check,
  Eye,
  EyeOff,
  Wallet,
  Coins,
  Lock,
  Sparkles,
  ShoppingBag,
} from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';
import { marketplaceApi } from '../services/api';
import type { Order } from '../types/marketplace';

export const OrdersPage: React.FC = () => {
  const { user, demoLogin } = useAuth();
  const { formatPrice } = useCart();

  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [revealed, setRevealed] = useState<Record<number, boolean>>({});
  const [copiedId, setCopiedId] = useState<number | null>(null);

  useEffect(() => {
    if (!user) {
      setLoading(false);
      return;
    }
    let active = true;
    setLoading(true);
    marketplaceApi
      .getOrders()
      .then((res) => {
        if (active) setOrders(res.orders || []);
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
  }, [user]);

  const handleCopy = (id: number, key: string) => {
    void navigator.clipboard.writeText(key);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  if (!user) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <div className="rounded-3xl bg-[#240a42] border-2 border-[#542294] p-10 space-y-5">
          <div className="w-14 h-14 rounded-2xl bg-[#ffd524]/15 border border-[#ffd524] flex items-center justify-center text-[#ffd524] mx-auto">
            <Lock className="w-7 h-7" />
          </div>
          <h1 className="text-2xl font-black text-white">Sign In Required to View Order History</h1>
          <p className="text-sm text-purple-300 max-w-md mx-auto">
            Your purchased digital keys are encrypted with AES-256-GCM and accessible exclusively to your authenticated buyer account.
          </p>
          <div className="flex items-center justify-center gap-3 flex-wrap">
            <button
              type="button"
              onClick={() => void demoLogin()}
              className="px-5 py-3 rounded-xl bg-[#00e676] text-[#16072c] font-black text-xs uppercase tracking-wider inline-flex items-center gap-2 cursor-pointer"
            >
              <Sparkles className="w-4 h-4" />
              <span>1-Click Instant Demo Login</span>
            </button>
            <Link
              to="/login?redirect=/orders"
              className="px-6 py-3 rounded-xl bg-[#ffd524] text-[#16072c] font-extrabold text-xs uppercase tracking-wider"
            >
              Sign In
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 py-8 space-y-8">
      {/* Buyer Account & Wallet Overview Banner */}
      <div className="rounded-3xl bg-gradient-to-r from-[#2e0d54] via-[#230940] to-[#1a0630] border border-[#542294] p-6 sm:p-8 flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div className="flex items-center gap-4">
          <div
            className="w-14 h-14 rounded-2xl flex items-center justify-center text-xl font-black text-[#16072c]"
            style={{ backgroundColor: user.avatar_color || '#ffd524' }}
          >
            {user.username.slice(0, 2).toUpperCase()}
          </div>
          <div>
            <span className="text-xs font-extrabold uppercase tracking-wider text-[#ffd524]">
              Verified Eneba Account
            </span>
            <h1 className="text-2xl font-black text-white">{user.username}</h1>
            <p className="text-xs text-purple-300">{user.email}</p>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="px-5 py-3.5 rounded-2xl bg-[#16052b] border border-[#491d7d]">
            <span className="text-[11px] text-purple-300 flex items-center gap-1.5">
              <Wallet className="w-3.5 h-3.5 text-[#00e676]" />
              Eneba Wallet Balance
            </span>
            <span className="text-xl font-black text-[#00e676] block mt-0.5">
              {formatPrice(user.wallet_balance)}
            </span>
          </div>

          <div className="px-5 py-3.5 rounded-2xl bg-[#16052b] border border-[#491d7d]">
            <span className="text-[11px] text-purple-300 flex items-center gap-1.5">
              <Coins className="w-3.5 h-3.5 text-[#ffd524]" />
              Total Cashback Earned
            </span>
            <span className="text-xl font-black text-[#ffd524] block mt-0.5">
              {formatPrice(user.cashback_earned)}
            </span>
          </div>
        </div>
      </div>

      {/* Orders List */}
      <div className="space-y-5">
        <div className="flex items-center justify-between">
          <h2 className="text-xl font-extrabold text-white flex items-center gap-2">
            <Package className="w-5 h-5 text-[#ffd524]" />
            Order History & Digital Key Library ({orders.length})
          </h2>
          <Link
            to="/store"
            className="text-xs font-bold text-[#ffd524] hover:underline"
          >
            + Buy More Games
          </Link>
        </div>

        {loading ? (
          <div className="h-48 rounded-3xl bg-[#240a42] animate-pulse border border-[#4c1e86]" />
        ) : orders.length === 0 ? (
          <div className="rounded-3xl bg-[#220940] border border-[#491d7d] p-12 text-center space-y-4">
            <ShoppingBag className="w-10 h-10 text-purple-300 mx-auto" />
            <h3 className="text-lg font-extrabold text-white">No Orders Yet</h3>
            <p className="text-xs text-purple-300 max-w-md mx-auto">
              When you purchase any game key or gift card, your instant decrypted activation codes will appear right here.
            </p>
            <Link
              to="/store"
              className="inline-block px-6 py-3 rounded-xl bg-[#ffd524] text-[#16072c] font-extrabold text-xs uppercase tracking-wider"
            >
              Browse Store Catalog
            </Link>
          </div>
        ) : (
          orders.map((order) => (
            <div
              key={order.id}
              className="rounded-3xl bg-[#220940] border border-[#4e1f88] overflow-hidden shadow-xl"
            >
              {/* Order Header */}
              <div className="px-6 py-4 bg-[#2b0c50] border-b border-[#461b7c] flex flex-wrap items-center justify-between gap-4 text-xs">
                <div className="flex items-center gap-4 flex-wrap">
                  <span className="font-mono font-extrabold text-sm text-[#ffd524]">
                    {order.order_number}
                  </span>
                  <span className="px-2.5 py-0.5 rounded-full bg-[#00e676]/20 text-[#00e676] font-extrabold">
                    {order.status}
                  </span>
                  <span className="text-purple-300">
                    {new Date(order.created_at).toLocaleString()}
                  </span>
                  <span className="text-purple-300">
                    via <strong className="text-white">{order.payment_provider}</strong>
                  </span>
                </div>
                <div className="font-extrabold text-sm text-white">
                  Total Paid: <span className="text-[#ffd524]">{formatPrice(order.total_amount)}</span>
                </div>
              </div>

              {/* Order Items & Keys */}
              <div className="p-6 space-y-4">
                {order.items.map((item) => {
                  const isShown = revealed[item.id] ?? false;
                  const masked = item.digital_key.replace(/[A-Z0-9]/g, '•');
                  return (
                    <div
                      key={item.id}
                      className="p-4 rounded-2xl bg-[#17052c] border border-[#3e176e] flex flex-col lg:flex-row lg:items-center justify-between gap-4"
                    >
                      <div className="flex items-center gap-3.5">
                        <img
                          src={item.product_cover}
                          alt={item.product_title}
                          className="w-12 h-16 object-cover rounded-lg border border-[#562396]"
                        />
                        <div>
                          <Link
                            to={`/product/${item.product_slug}`}
                            className="text-sm font-extrabold text-white hover:text-[#ffd524]"
                          >
                            {item.product_title}
                          </Link>
                          <div className="flex items-center gap-2 mt-1 text-xs text-purple-300">
                            <span className="px-2 py-0.5 rounded bg-[#32115c] text-purple-100 text-[10px] font-bold">
                              {item.platform}
                            </span>
                            <span className="px-2 py-0.5 rounded bg-[#00e676]/15 text-[#00e676] text-[10px] font-bold">
                              {item.region}
                            </span>
                            <span>Seller: {item.seller_name}</span>
                          </div>
                        </div>
                      </div>

                      <div className="flex flex-col sm:flex-row sm:items-center gap-3 bg-[#240a42] border border-[#542294] px-4 py-2.5 rounded-xl">
                        <div className="flex items-center gap-2">
                          <KeyRound className="w-4 h-4 text-[#ffd524] shrink-0" />
                          <code className="font-mono font-extrabold text-sm text-[#ffd524] tracking-wider">
                            {isShown ? item.digital_key : masked}
                          </code>
                        </div>
                        <div className="flex items-center gap-2">
                          <button
                            type="button"
                            onClick={() =>
                              setRevealed((prev) => ({ ...prev, [item.id]: !isShown }))
                            }
                            className="px-2.5 py-1.5 rounded-lg bg-[#381366] hover:bg-[#491985] text-xs font-bold text-white inline-flex items-center gap-1 cursor-pointer"
                          >
                            {isShown ? <EyeOff className="w-3.5 h-3.5" /> : <Eye className="w-3.5 h-3.5" />}
                            <span>{isShown ? 'Hide' : 'Reveal Key'}</span>
                          </button>
                          <button
                            type="button"
                            onClick={() => handleCopy(item.id, item.digital_key)}
                            className="px-3 py-1.5 rounded-lg bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] text-xs font-extrabold inline-flex items-center gap-1 cursor-pointer"
                          >
                            {copiedId === item.id ? (
                              <>
                                <Check className="w-3.5 h-3.5" />
                                <span>Copied</span>
                              </>
                            ) : (
                              <>
                                <Copy className="w-3.5 h-3.5" />
                                <span>Copy</span>
                              </>
                            )}
                          </button>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};
