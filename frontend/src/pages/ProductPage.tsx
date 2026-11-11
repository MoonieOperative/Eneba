import React, { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import {
  ShieldCheck,
  Zap,
  Globe,
  Star,
  Heart,
  ShoppingCart,
  CheckCircle2,
  Coins,
  Award,
  Cpu,
  KeyRound,
  ChevronRight,
  Users,
  Sparkles,
} from 'lucide-react';
import { marketplaceApi } from '../services/api';
import type { Offer, Product } from '../types/marketplace';
import { useCart } from '../context/CartContext';

export const ProductPage: React.FC = () => {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { addToCart, toggleWishlist, isWishlisted, formatPrice } = useCart();

  const [product, setProduct] = useState<Product | null>(null);
  const [selectedOffer, setSelectedOffer] = useState<Offer | null>(null);
  const [activeImage, setActiveImage] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!slug) return;
    let active = true;
    setLoading(true);
    marketplaceApi
      .getProductBySlug(slug)
      .then((res) => {
        if (active) {
          setProduct(res);
          setSelectedOffer(res.offers && res.offers.length > 0 ? res.offers[0] : null);
          setActiveImage(res.banner_image || res.cover_image);
          setError(null);
        }
      })
      .catch((err: Error) => {
        if (active) setError(err.message);
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [slug]);

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-12 space-y-8">
        <div className="h-96 rounded-3xl bg-[#260c47] animate-pulse border border-[#491d7d]" />
        <div className="h-64 rounded-3xl bg-[#260c47] animate-pulse border border-[#491d7d]" />
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <div className="p-8 rounded-3xl bg-[#260c47] border border-[#4c1e86] space-y-4">
          <h2 className="text-2xl font-extrabold text-white">Product Not Found</h2>
          <p className="text-sm text-purple-300">{error || 'The requested game key could not be found.'}</p>
          <Link
            to="/store"
            className="inline-block px-6 py-3 rounded-xl bg-[#ffd524] text-[#16072c] font-extrabold text-xs uppercase tracking-wider"
          >
            Back to Store Catalog
          </Link>
        </div>
      </div>
    );
  }

  const offers = product.offers || [];
  const currentOffer = selectedOffer || offers[0];
  const wishlisted = isWishlisted(product.id);

  const handleBuyNow = (offerToBuy?: Offer) => {
    addToCart(product, offerToBuy || currentOffer, false);
    navigate('/checkout');
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-10">
      {/* Breadcrumb Navigation */}
      <nav className="flex items-center gap-2 text-xs text-purple-300 overflow-x-auto no-scrollbar">
        <Link to="/" className="hover:text-[#ffd524]">Home</Link>
        <ChevronRight className="w-3.5 h-3.5 text-purple-500 shrink-0" />
        <Link to="/store" className="hover:text-[#ffd524]">Store</Link>
        <ChevronRight className="w-3.5 h-3.5 text-purple-500 shrink-0" />
        <Link to={`/store?platform=${encodeURIComponent(product.platform)}`} className="hover:text-[#ffd524] shrink-0">
          {product.platform}
        </Link>
        <ChevronRight className="w-3.5 h-3.5 text-purple-500 shrink-0" />
        <span className="text-white font-bold truncate">{product.title}</span>
      </nav>

      {/* Top Product Showcase Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left Cover + Gallery (3 cols) */}
        <div className="lg:col-span-3 space-y-4">
          <div className="relative rounded-2xl overflow-hidden border-2 border-[#542294] shadow-2xl bg-[#190630]">
            <img
              src={product.cover_image}
              alt={product.title}
              className="w-full aspect-[3/4] object-cover"
            />
            {product.discount_percent > 0 && (
              <span className="absolute top-3 left-3 px-2.5 py-1 rounded-lg bg-[#ff4081] text-white font-extrabold text-xs shadow-lg">
                SAVE {product.discount_percent}%
              </span>
            )}
          </div>

          <button
            type="button"
            onClick={() => toggleWishlist(product.id)}
            className={`w-full py-2.5 px-4 rounded-xl border font-bold text-xs flex items-center justify-center gap-2 transition-colors cursor-pointer ${
              wishlisted
                ? 'bg-[#ff4081]/20 border-[#ff4081] text-white'
                : 'bg-[#250b45] border-[#4b1e82] text-purple-200 hover:text-white hover:border-[#ffd524]'
            }`}
          >
            <Heart className={`w-4 h-4 ${wishlisted ? 'fill-[#ff4081] text-[#ff4081]' : ''}`} />
            <span>{wishlisted ? 'Saved in Wishlist' : `Add to Wishlist (${product.likes_count.toLocaleString()})`}</span>
          </button>
        </div>

        {/* Center Product Specs & Verification Badges (5 cols) */}
        <div className="lg:col-span-5 space-y-5">
          <div>
            <div className="flex items-center gap-2 flex-wrap mb-2">
              <span className="px-2.5 py-0.5 rounded-md bg-[#3c166e] text-[#ffd524] font-extrabold text-xs">
                {product.product_type}
              </span>
              <span className="px-2.5 py-0.5 rounded-md bg-[#2c0e52] text-purple-200 font-bold text-xs">
                {product.genre}
              </span>
              <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-md bg-[#ffd524]/15 text-[#ffd524] font-bold text-xs">
                <Star className="w-3.5 h-3.5 fill-[#ffd524]" />
                {product.rating.toFixed(2)} ({product.review_count.toLocaleString()} reviews)
              </span>
            </div>

            <h1 className="text-2xl sm:text-3xl font-black text-white leading-tight">
              {product.title}
            </h1>
            <p className="text-sm text-purple-200 mt-1">{product.subtitle}</p>
          </div>

          {/* Eneba Key Compatibility & Activation Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 p-4 rounded-2xl bg-[#220940] border border-[#4a1d80]">
            <div className="flex items-start gap-3">
              <div className="w-9 h-9 rounded-xl bg-[#00e676]/15 border border-[#00e676]/30 flex items-center justify-center text-[#00e676] shrink-0">
                <Globe className="w-5 h-5" />
              </div>
              <div>
                <span className="text-[11px] uppercase tracking-wider text-purple-300 font-bold block">
                  Region Lock Check
                </span>
                <span className="text-xs font-extrabold text-[#00e676]">
                  ✓ Activates in {product.region}
                </span>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="w-9 h-9 rounded-xl bg-[#ffd524]/15 border border-[#ffd524]/30 flex items-center justify-center text-[#ffd524] shrink-0">
                <KeyRound className="w-5 h-5" />
              </div>
              <div>
                <span className="text-[11px] uppercase tracking-wider text-purple-300 font-bold block">
                  Platform / DRM
                </span>
                <span className="text-xs font-extrabold text-white">
                  {product.platform} Official Key
                </span>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="w-9 h-9 rounded-xl bg-[#00b0ff]/15 border border-[#00b0ff]/30 flex items-center justify-center text-[#00b0ff] shrink-0">
                <Zap className="w-5 h-5" />
              </div>
              <div>
                <span className="text-[11px] uppercase tracking-wider text-purple-300 font-bold block">
                  Delivery Speed
                </span>
                <span className="text-xs font-extrabold text-white">
                  Instant Automated Dispatch
                </span>
              </div>
            </div>

            <div className="flex items-start gap-3">
              <div className="w-9 h-9 rounded-xl bg-[#ff4081]/15 border border-[#ff4081]/30 flex items-center justify-center text-[#ff4081] shrink-0">
                <ShieldCheck className="w-5 h-5" />
              </div>
              <div>
                <span className="text-[11px] uppercase tracking-wider text-purple-300 font-bold block">
                  Key Security
                </span>
                <span className="text-xs font-extrabold text-white">
                  AES-256-GCM Encrypted
                </span>
              </div>
            </div>
          </div>

          {/* Media Preview / Screenshots */}
          <div className="space-y-2.5">
            <div className="rounded-2xl overflow-hidden border border-[#4a1d80] bg-[#16052b] aspect-video">
              <img
                src={activeImage || product.banner_image || product.cover_image}
                alt={product.title}
                className="w-full h-full object-cover"
              />
            </div>
            {product.screenshots && product.screenshots.length > 0 && (
              <div className="flex items-center gap-2 overflow-x-auto pb-1">
                {[product.banner_image, ...product.screenshots].filter(Boolean).map((imgUrl, idx) => (
                  <button
                    key={idx}
                    type="button"
                    onClick={() => setActiveImage(imgUrl)}
                    className={`w-24 aspect-video rounded-lg overflow-hidden border-2 shrink-0 cursor-pointer ${
                      activeImage === imgUrl ? 'border-[#ffd524]' : 'border-transparent opacity-70 hover:opacity-100'
                    }`}
                  >
                    <img src={imgUrl} alt="" className="w-full h-full object-cover" />
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Right Primary Buy Box (4 cols) */}
        <div className="lg:col-span-4 rounded-3xl bg-gradient-to-b from-[#2f0e59] to-[#21093e] border-2 border-[#ffd524]/70 p-6 shadow-2xl space-y-5 sticky top-28">
          {currentOffer && (
            <>
              <div className="flex items-center justify-between gap-2 pb-3 border-b border-[#491d7d]">
                <div>
                  <span className="text-[10px] font-extrabold uppercase tracking-wider text-[#ffd524] block">
                    Selected Merchant Offer
                  </span>
                  <div className="flex items-center gap-1.5 mt-0.5">
                    <span className="text-sm font-extrabold text-white">
                      {currentOffer.seller.name}
                    </span>
                    <CheckCircle2 className="w-4 h-4 text-[#00e676]" />
                  </div>
                </div>
                <div className="text-right">
                  <span className="px-2 py-0.5 rounded bg-[#00e676]/20 text-[#00e676] font-extrabold text-xs">
                    {currentOffer.seller.rating_percent}% Positive
                  </span>
                  <span className="block text-[10px] text-purple-300 mt-0.5">
                    {currentOffer.seller.total_sales.toLocaleString()} sales
                  </span>
                </div>
              </div>

              {/* Price Display */}
              <div className="space-y-1.5">
                <div className="flex items-baseline justify-between">
                  <div>
                    <span className="text-3xl sm:text-4xl font-black text-white">
                      {formatPrice(currentOffer.price)}
                    </span>
                    {product.base_price > currentOffer.price && (
                      <span className="text-sm text-purple-300 line-through ml-2.5">
                        {formatPrice(product.base_price)}
                      </span>
                    )}
                  </div>
                  <span className="text-xs font-extrabold text-[#00e676]">
                    ● In Stock ({currentOffer.stock_count} keys)
                  </span>
                </div>

                <div className="flex items-center gap-2 px-3 py-2 rounded-xl bg-[#00e676]/15 border border-[#00e676]/30 text-xs font-bold text-[#00e676]">
                  <Coins className="w-4 h-4 shrink-0" />
                  <span>
                    Earn +{formatPrice(currentOffer.cashback_amount)} ({currentOffer.cashback_percent}%) Eneba Cashback
                  </span>
                </div>
              </div>

              {/* Primary CTAs */}
              <div className="space-y-2.5">
                <button
                  type="button"
                  onClick={() => handleBuyNow(currentOffer)}
                  className="w-full py-3.5 px-6 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-black text-sm uppercase tracking-wider flex items-center justify-center gap-2 shadow-xl shadow-[#ffd524]/25 transition-transform active:scale-98 cursor-pointer"
                >
                  <Zap className="w-4 h-4 fill-current" />
                  <span>Buy Now — Instant Key</span>
                </button>

                <button
                  type="button"
                  onClick={() => addToCart(product, currentOffer, true)}
                  className="w-full py-3 px-6 rounded-xl bg-[#3c166e] hover:bg-[#4d1d8c] border border-[#682cb3] text-white font-extrabold text-xs uppercase tracking-wider flex items-center justify-center gap-2 transition-colors cursor-pointer"
                >
                  <ShoppingCart className="w-4 h-4 text-[#ffd524]" />
                  <span>Add to Cart</span>
                </button>
              </div>

              {/* Jump to Multi-Seller Comparison Note */}
              <div className="pt-3 border-t border-[#491d7d] flex items-center justify-between text-xs text-purple-200">
                <span className="flex items-center gap-1.5">
                  <Users className="w-4 h-4 text-[#ffd524]" />
                  <strong>{offers.length} sellers</strong> competing on this product
                </span>
                <a href="#seller-offers" className="text-[#ffd524] font-bold hover:underline">
                  Compare All ↓
                </a>
              </div>
            </>
          )}
        </div>
      </div>

      {/* Multi-Seller Offers Comparison Section */}
      <section id="seller-offers" className="rounded-3xl bg-[#220940] border-2 border-[#4d1f85] p-6 sm:p-8 space-y-5">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-4 border-b border-[#3b156b]">
          <div>
            <span className="text-xs font-extrabold uppercase tracking-wider text-[#ffd524] flex items-center gap-1.5">
              <Award className="w-4 h-4" />
              Marketplace Seller Competition
            </span>
            <h2 className="text-xl font-extrabold text-white mt-0.5">
              Compare Offers from {offers.length} Verified Sellers
            </h2>
          </div>
          <span className="text-xs text-purple-300">
            Click any merchant below to select their offer or add directly to your cart
          </span>
        </div>

        <div className="space-y-3">
          {offers.map((offer, index) => {
            const isSelected = currentOffer?.id === offer.id;
            return (
              <div
                key={offer.id}
                onClick={() => setSelectedOffer(offer)}
                className={`p-4 rounded-2xl border-2 transition-all flex flex-col md:flex-row md:items-center justify-between gap-4 cursor-pointer ${
                  isSelected
                    ? 'bg-[#31115b] border-[#ffd524] shadow-lg shadow-[#ffd524]/10'
                    : 'bg-[#1a0633] border-[#3d176e] hover:border-[#682cb3]'
                }`}
              >
                {/* Seller Identity & Metrics */}
                <div className="flex items-center gap-4">
                  <div className="w-11 h-11 rounded-xl bg-[#2a0c4e] border border-[#562396] flex items-center justify-center font-black text-sm text-[#ffd524] shrink-0">
                    {offer.seller.country_code}
                  </div>
                  <div>
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-sm font-extrabold text-white">{offer.seller.name}</span>
                      {offer.seller.verified && (
                        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded bg-[#00e676]/15 text-[#00e676] text-[10px] font-extrabold">
                          <CheckCircle2 className="w-3 h-3" /> {offer.seller.badge}
                        </span>
                      )}
                      {index === 0 && (
                        <span className="px-2 py-0.5 rounded bg-[#ffd524] text-[#16072c] text-[10px] font-black uppercase">
                          Lowest Price
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-4 mt-1 text-xs text-purple-300">
                      <span className="text-[#ffd524] font-bold">
                        ★ {offer.seller.rating_percent}% positive feedback
                      </span>
                      <span>• {offer.seller.total_sales.toLocaleString()} completed orders</span>
                      <span className="hidden sm:inline">• Joined {offer.seller.joined_year}</span>
                    </div>
                  </div>
                </div>

                {/* Delivery, Cashback, Price & Action Buttons */}
                <div className="flex items-center justify-between md:justify-end gap-5 pt-3 md:pt-0 border-t md:border-t-0 border-[#34125e]">
                  <div className="text-left md:text-right">
                    <div className="text-xl font-black text-white">{formatPrice(offer.price)}</div>
                    <div className="text-[11px] font-bold text-[#00e676] flex items-center md:justify-end gap-1">
                      <Sparkles className="w-3 h-3" />
                      +{formatPrice(offer.cashback_amount)} ({offer.cashback_percent}% cashback)
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <button
                      type="button"
                      onClick={(e) => {
                        e.stopPropagation();
                        addToCart(product, offer, true);
                      }}
                      className="px-4 py-2.5 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs uppercase tracking-wide flex items-center gap-1.5 cursor-pointer"
                    >
                      <ShoppingCart className="w-3.5 h-3.5" />
                      <span>Add to Cart</span>
                    </button>
                    <button
                      type="button"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleBuyNow(offer);
                      }}
                      className="px-3.5 py-2.5 rounded-xl bg-[#3c166e] hover:bg-[#4e1d8c] text-white font-bold text-xs cursor-pointer"
                    >
                      Buy Now
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </section>

      {/* Product Description, System Requirements & Activation Guide */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div className="lg:col-span-7 rounded-3xl bg-[#220940] border border-[#481d7d] p-6 sm:p-8 space-y-5">
          <h3 className="text-lg font-extrabold text-white">About {product.title}</h3>
          <p className="text-sm text-purple-200 leading-relaxed">{product.description}</p>

          <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 pt-4 border-t border-[#3b156b] text-xs">
            <div>
              <span className="text-purple-400 block">Developer</span>
              <strong className="text-white">{product.developer}</strong>
            </div>
            <div>
              <span className="text-purple-400 block">Publisher</span>
              <strong className="text-white">{product.publisher}</strong>
            </div>
            <div>
              <span className="text-purple-400 block">Release Date</span>
              <strong className="text-white">{product.release_date}</strong>
            </div>
          </div>

          {product.activation_steps && product.activation_steps.length > 0 && (
            <div className="pt-5 border-t border-[#3b156b] space-y-3">
              <h4 className="text-sm font-extrabold text-[#ffd524] uppercase tracking-wider">
                How to Activate Your {product.platform} Key
              </h4>
              <ol className="space-y-2 text-xs text-purple-200 list-decimal list-inside">
                {product.activation_steps.map((step, idx) => (
                  <li key={idx} className="leading-relaxed">{step}</li>
                ))}
              </ol>
            </div>
          )}
        </div>

        <div className="lg:col-span-5 rounded-3xl bg-[#220940] border border-[#481d7d] p-6 sm:p-8 space-y-4">
          <div className="flex items-center gap-2">
            <Cpu className="w-5 h-5 text-[#ffd524]" />
            <h3 className="text-lg font-extrabold text-white">System Requirements</h3>
          </div>
          <div className="space-y-3 text-xs">
            <div className="p-3 rounded-xl bg-[#18062e] border border-[#3d176e]">
              <span className="text-purple-400 block font-semibold">Operating System</span>
              <span className="text-white font-bold">{product.system_requirements.os}</span>
            </div>
            <div className="p-3 rounded-xl bg-[#18062e] border border-[#3d176e]">
              <span className="text-purple-400 block font-semibold">Processor (CPU)</span>
              <span className="text-white font-bold">{product.system_requirements.processor}</span>
            </div>
            <div className="p-3 rounded-xl bg-[#18062e] border border-[#3d176e]">
              <span className="text-purple-400 block font-semibold">Memory (RAM)</span>
              <span className="text-white font-bold">{product.system_requirements.memory}</span>
            </div>
            <div className="p-3 rounded-xl bg-[#18062e] border border-[#3d176e]">
              <span className="text-purple-400 block font-semibold">Graphics (GPU)</span>
              <span className="text-white font-bold">{product.system_requirements.graphics}</span>
            </div>
            <div className="p-3 rounded-xl bg-[#18062e] border border-[#3d176e]">
              <span className="text-purple-400 block font-semibold">Storage</span>
              <span className="text-white font-bold">{product.system_requirements.storage}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
