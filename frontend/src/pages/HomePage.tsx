import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import {
  ArrowRight,
  Flame,
  Sparkles,
  Server,
  ShieldCheck,
  Zap,
  Lock,
  Users,
  Coins,
  Monitor,
  Gamepad2,
  Trophy,
  ShoppingCart,
} from 'lucide-react';
import { marketplaceApi } from '../services/api';
import type { HomepageBlock } from '../types/marketplace';
import { ProductCard } from '../components/ProductCard';
import { useCart } from '../context/CartContext';

export const HomePage: React.FC = () => {
  const [blocks, setBlocks] = useState<HomepageBlock[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    setLoading(true);
    marketplaceApi
      .getHomepageBlocks()
      .then((res) => {
        if (active) {
          setBlocks(res.blocks || []);
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
  }, []);

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-12 space-y-10">
        <div className="h-80 rounded-3xl bg-[#270c47] animate-pulse border border-[#491d7d]" />
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-5">
          {Array.from({ length: 5 }).map((_, idx) => (
            <div
              key={idx}
              className="aspect-[3/4] rounded-2xl bg-[#270c47] animate-pulse border border-[#491d7d]"
            />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-16 text-center">
        <div className="p-8 rounded-3xl bg-[#280d4a] border border-red-500/40 space-y-4">
          <h2 className="text-xl font-extrabold text-white">Unable to Load Dynamic Homepage Blocks</h2>
          <p className="text-sm text-purple-200">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8 space-y-14">
      {blocks.map((block) => (
        <DynamicHomepageBlock key={block.id || block.slug} block={block} />
      ))}
    </div>
  );
};

const DynamicHomepageBlock: React.FC<{ block: HomepageBlock }> = ({ block }) => {
  const { addToCart, formatPrice } = useCart();

  // 1. Hero Banner Block
  if (block.block_type === 'hero_banner') {
    const heroProducts = block.products || [];
    const primaryProduct = heroProducts[0];
    const sideProducts = heroProducts.slice(1, 3);
    const promoCodeHint = (block.config?.promo_code_hint as string) || '';
    const promoDiscount = (block.config?.promo_discount as string) || '';

    return (
      <section className="relative rounded-3xl overflow-hidden bg-gradient-to-br from-[#381069] via-[#270a4a] to-[#17052e] border-2 border-[#5c259e] shadow-2xl p-6 sm:p-10">
        {/* Background Glow */}
        <div
          className="absolute -top-24 -right-24 w-96 h-96 rounded-full blur-3xl opacity-25 pointer-events-none"
          style={{ backgroundColor: block.accent_color || '#ffd524' }}
        />

        <div className="relative z-10 grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
          {/* Left Promotional Copy (from DB block metadata) */}
          <div className="lg:col-span-6 space-y-5">
            <div className="flex items-center gap-2.5 flex-wrap">
              {block.badge_text && (
                <span
                  className="px-3 py-1 rounded-full text-xs font-extrabold uppercase tracking-wider text-[#16072c]"
                  style={{ backgroundColor: block.accent_color || '#ffd524' }}
                >
                  {block.badge_text}
                </span>
              )}
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-[#150529]/80 border border-[#52218c] text-[11px] font-mono text-purple-200">
                <Server className="w-3 h-3 text-[#00e676]" />
                {block.service_name}
              </span>
            </div>

            <h1 className="text-3xl sm:text-5xl font-black tracking-tight text-white leading-[1.1]">
              {block.title}
            </h1>

            <p className="text-sm sm:text-base text-purple-200 leading-relaxed max-w-xl">
              {block.subtitle}
            </p>

            {promoCodeHint && (
              <div className="inline-flex items-center gap-3 px-4 py-2.5 rounded-2xl bg-[#180630]/90 border border-[#ffd524]/40">
                <Sparkles className="w-4 h-4 text-[#ffd524]" />
                <span className="text-xs text-purple-100">
                  Promo Code:{' '}
                  <code className="font-mono font-extrabold text-[#ffd524] bg-[#2e1054] px-2 py-0.5 rounded">
                    {promoCodeHint}
                  </code>
                </span>
                {promoDiscount && (
                  <span className="text-[11px] font-extrabold text-[#00e676]">{promoDiscount}</span>
                )}
              </div>
            )}

            <div className="pt-1 flex items-center gap-4 flex-wrap">
              {block.cta_link && (
                <Link
                  to={block.cta_link}
                  className="px-6 py-3.5 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-sm uppercase tracking-wider inline-flex items-center gap-2 shadow-xl shadow-[#ffd524]/20 transition-transform active:scale-95"
                >
                  <span>{block.cta_text || 'Shop Now'}</span>
                  <ArrowRight className="w-4 h-4" />
                </Link>
              )}
              <Link
                to="/store"
                className="px-5 py-3.5 rounded-xl bg-[#2a0c4e] hover:bg-[#381266] border border-[#5c259e] text-white font-bold text-sm transition-colors"
              >
                Browse All Games
              </Link>
            </div>
          </div>

          {/* Right Featured Spotlight Cards (dynamically resolved by HeroPromotionService) */}
          <div className="lg:col-span-6 grid grid-cols-1 sm:grid-cols-2 gap-4">
            {primaryProduct && (
              <div className="sm:col-span-2 rounded-2xl bg-[#1b0733]/90 border border-[#5c259e] overflow-hidden flex flex-col sm:flex-row group">
                <Link to={`/product/${primaryProduct.slug}`} className="sm:w-48 shrink-0 overflow-hidden">
                  <img
                    src={primaryProduct.banner_image || primaryProduct.cover_image}
                    alt={primaryProduct.title}
                    className="w-full h-44 sm:h-full object-cover group-hover:scale-105 transition-transform duration-300"
                  />
                </Link>
                <div className="p-4 flex-1 flex flex-col justify-between gap-3">
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-[10px] font-extrabold px-2 py-0.5 rounded bg-[#ff4081] text-white">
                        -{primaryProduct.discount_percent}% FLASH
                      </span>
                      <span className="text-[10px] font-extrabold px-2 py-0.5 rounded bg-[#00e676]/20 text-[#00e676]">
                        {primaryProduct.region}
                      </span>
                    </div>
                    <Link
                      to={`/product/${primaryProduct.slug}`}
                      className="text-base font-extrabold text-white hover:text-[#ffd524] transition-colors line-clamp-1"
                    >
                      {primaryProduct.title}
                    </Link>
                    <p className="text-xs text-purple-300 line-clamp-2 mt-1">{primaryProduct.subtitle}</p>
                  </div>
                  <div className="flex items-center justify-between pt-2 border-t border-[#3b156b]">
                    <div>
                      <span className="text-[11px] text-purple-400 line-through mr-1.5">
                        {formatPrice(primaryProduct.base_price)}
                      </span>
                      <span className="text-xl font-black text-[#ffd524]">
                        {formatPrice(primaryProduct.lowest_price)}
                      </span>
                    </div>
                    <button
                      type="button"
                      onClick={() => addToCart(primaryProduct)}
                      className="px-3.5 py-2 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs inline-flex items-center gap-1.5 cursor-pointer"
                    >
                      <ShoppingCart className="w-3.5 h-3.5" />
                      Add to Cart
                    </button>
                  </div>
                </div>
              </div>
            )}

            {sideProducts.map((prod) => (
              <Link
                key={prod.id}
                to={`/product/${prod.slug}`}
                className="rounded-2xl bg-[#1b0733]/90 hover:bg-[#280c4a] border border-[#4e1f88] p-3 flex items-center gap-3 transition-colors group"
              >
                <img
                  src={prod.cover_image}
                  alt={prod.title}
                  className="w-14 h-20 object-cover rounded-xl shrink-0 border border-[#5a259c]"
                />
                <div className="min-w-0 flex-1">
                  <span className="text-[10px] font-extrabold text-[#00e676]">{prod.platform}</span>
                  <h3 className="text-xs font-extrabold text-white group-hover:text-[#ffd524] truncate">
                    {prod.title}
                  </h3>
                  <div className="mt-1 flex items-center gap-2">
                    <span className="text-sm font-black text-[#ffd524]">
                      {formatPrice(prod.lowest_price)}
                    </span>
                    {prod.discount_percent > 0 && (
                      <span className="text-[10px] font-extrabold px-1.5 py-0.5 rounded bg-[#ff4081] text-white">
                        -{prod.discount_percent}%
                      </span>
                    )}
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>
    );
  }

  // 2. Category Grid Block
  if (block.block_type === 'category_grid') {
    const items = block.items || [];
    return (
      <section className="space-y-5">
        <BlockSectionHeader block={block} />
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {items.map((item, idx) => (
            <Link
              key={idx}
              to={item.link || '/store'}
              className="group relative rounded-2xl bg-[#250b45] hover:bg-[#31115b] border border-[#491d7d] hover:border-[#ffd524]/60 p-5 flex items-start justify-between gap-4 transition-all hover:-translate-y-0.5 shadow-lg"
            >
              <div className="space-y-1.5">
                {item.badge && (
                  <span
                    className="inline-block text-[10px] font-extrabold px-2 py-0.5 rounded text-[#16072c]"
                    style={{ backgroundColor: item.accent_color || '#ffd524' }}
                  >
                    {item.badge}
                  </span>
                )}
                <h3 className="text-base font-extrabold text-white group-hover:text-[#ffd524] transition-colors">
                  {item.title}
                </h3>
                <p className="text-xs text-purple-300">{item.subtitle}</p>
                {item.metric && (
                  <p className="text-[11px] font-bold text-[#00e676] pt-1">{item.metric}</p>
                )}
              </div>
              <div
                className="w-11 h-11 rounded-xl flex items-center justify-center shrink-0 border border-white/10"
                style={{ backgroundColor: `${item.accent_color || '#ffd524'}22`, color: item.accent_color || '#ffd524' }}
              >
                {renderDynamicIcon(item.icon)}
              </div>
            </Link>
          ))}
        </div>
      </section>
    );
  }

  // 3. Value Props / Trust Shield Block
  if (block.block_type === 'value_props') {
    const items = block.items || [];
    return (
      <section className="rounded-3xl bg-gradient-to-r from-[#220940] via-[#2b0d4e] to-[#220940] border border-[#4e1f88] p-6 sm:p-8 space-y-6">
        <BlockSectionHeader block={block} />
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          {items.map((item, idx) => (
            <div
              key={idx}
              className="rounded-2xl bg-[#18062e]/90 border border-[#421973] p-5 space-y-2.5"
            >
              <div className="flex items-center justify-between">
                <div
                  className="w-10 h-10 rounded-xl flex items-center justify-center"
                  style={{ backgroundColor: `${item.accent_color || '#ffd524'}20`, color: item.accent_color || '#ffd524' }}
                >
                  {renderDynamicIcon(item.icon)}
                </div>
                {item.badge && (
                  <span className="text-[10px] font-mono font-bold px-2 py-0.5 rounded bg-[#2d0f52] text-[#ffd524]">
                    {item.badge}
                  </span>
                )}
              </div>
              <h3 className="text-sm font-extrabold text-white">{item.title}</h3>
              <p className="text-xs text-purple-300 leading-relaxed">{item.subtitle}</p>
            </div>
          ))}
        </div>
      </section>
    );
  }

  // 4. Flash Deals or Standard Product Grid Block
  const products = block.products || [];
  return (
    <section className="space-y-5">
      <BlockSectionHeader block={block} />
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 sm:gap-5">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </section>
  );
};

const BlockSectionHeader: React.FC<{ block: HomepageBlock }> = ({ block }) => {
  return (
    <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-3 pb-2 border-b border-[#3b156b]">
      <div className="space-y-1">
        <div className="flex items-center gap-2 flex-wrap">
          {block.badge_text && (
            <span
              className="px-2.5 py-0.5 rounded-md text-[10px] font-extrabold uppercase tracking-wider text-[#16072c] inline-flex items-center gap-1"
              style={{ backgroundColor: block.accent_color || '#ffd524' }}
            >
              <Flame className="w-3 h-3 fill-current" />
              {block.badge_text}
            </span>
          )}
          <span className="inline-flex items-center gap-1 text-[11px] font-mono text-purple-300 bg-[#230942] px-2 py-0.5 rounded border border-[#431978]">
            <Server className="w-3 h-3 text-[#00e676]" />
            {block.service_name}
          </span>
        </div>
        <h2 className="text-xl sm:text-2xl font-extrabold text-white tracking-tight">
          {block.title}
        </h2>
        {block.subtitle && <p className="text-xs text-purple-300">{block.subtitle}</p>}
      </div>

      {block.cta_link && (
        <Link
          to={block.cta_link}
          className="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#34125e] hover:bg-[#ffd524] text-white hover:text-[#16072c] text-xs font-extrabold transition-colors shrink-0 self-start sm:self-auto"
        >
          <span>{block.cta_text || 'View All'}</span>
          <ArrowRight className="w-3.5 h-3.5" />
        </Link>
      )}
    </div>
  );
};

function renderDynamicIcon(iconName?: string) {
  switch (iconName) {
    case 'lock':
      return <Lock className="w-5 h-5" />;
    case 'users':
      return <Users className="w-5 h-5" />;
    case 'zap':
      return <Zap className="w-5 h-5" />;
    case 'coins':
      return <Coins className="w-5 h-5" />;
    case 'monitor':
      return <Monitor className="w-5 h-5" />;
    case 'gamepad-2':
      return <Gamepad2 className="w-5 h-5" />;
    case 'trophy':
      return <Trophy className="w-5 h-5" />;
    case 'sparkles':
      return <Sparkles className="w-5 h-5" />;
    default:
      return <ShieldCheck className="w-5 h-5" />;
  }
}
