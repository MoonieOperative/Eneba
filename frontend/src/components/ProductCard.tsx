import React from 'react';
import { Link } from 'react-router-dom';
import { Heart, ShoppingCart, Star, Coins, Layers } from 'lucide-react';
import type { Product } from '../types/marketplace';
import { useCart } from '../context/CartContext';

interface ProductCardProps {
  product: Product;
}

export const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const { addToCart, toggleWishlist, isWishlisted, formatPrice } = useCart();
  const wishlisted = isWishlisted(product.id);

  const regionColorClass =
    product.region === 'GLOBAL'
      ? 'text-[#00e676] bg-[#00e676]/15 border-[#00e676]/30'
      : product.region === 'EUROPE'
      ? 'text-[#00b0ff] bg-[#00b0ff]/15 border-[#00b0ff]/30'
      : 'text-[#ffd524] bg-[#ffd524]/15 border-[#ffd524]/30';

  return (
    <div className="group relative rounded-2xl bg-[#260c47] hover:bg-[#31115b] border border-[#471c7c] hover:border-[#ffd524]/70 transition-all duration-200 flex flex-col overflow-hidden shadow-lg hover:shadow-2xl hover:shadow-[#ffd524]/10 hover:-translate-y-1">
      {/* Cover Image Container */}
      <div className="relative aspect-[3/4] w-full overflow-hidden bg-[#18062e]">
        <Link to={`/product/${product.slug}`} className="block w-full h-full">
          <img
            src={product.cover_image}
            alt={product.title}
            loading="lazy"
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
        </Link>

        {/* Top Left Badges: Discount & Cashback */}
        <div className="absolute top-2.5 left-2.5 flex flex-col gap-1.5 pointer-events-none">
          {product.discount_percent > 0 && (
            <span className="px-2 py-0.5 rounded-md bg-[#ff4081] text-white font-extrabold text-xs shadow-md">
              -{product.discount_percent}%
            </span>
          )}
          {product.cashback_amount > 0 && (
            <span className="px-2 py-0.5 rounded-md bg-[#16072c]/90 backdrop-blur-xs border border-[#00e676]/50 text-[#00e676] font-bold text-[10px] flex items-center gap-1 shadow-md">
              <Coins className="w-3 h-3" />
              +{formatPrice(product.cashback_amount)}
            </span>
          )}
        </div>

        {/* Top Right Wishlist Heart Button */}
        <button
          type="button"
          onClick={(e) => {
            e.preventDefault();
            toggleWishlist(product.id);
          }}
          aria-label="Toggle wishlist"
          className={`absolute top-2.5 right-2.5 w-8 h-8 rounded-full flex items-center justify-center backdrop-blur-md border transition-colors cursor-pointer ${
            wishlisted
              ? 'bg-[#ff4081] border-[#ff4081] text-white'
              : 'bg-[#16072c]/75 border-white/20 text-white hover:bg-[#ff4081] hover:border-[#ff4081]'
          }`}
        >
          <Heart className={`w-4 h-4 ${wishlisted ? 'fill-white' : ''}`} />
        </button>

        {/* Bottom Platform & Offers Count Overlay Strip */}
        <div className="absolute bottom-0 inset-x-0 px-3 py-1.5 bg-gradient-to-t from-[#140526] via-[#140526]/90 to-transparent flex items-center justify-between text-[11px]">
          <span className="font-bold text-purple-100 truncate">{product.platform}</span>
          <span className="inline-flex items-center gap-1 text-purple-300 font-semibold shrink-0">
            <Layers className="w-3 h-3 text-[#ffd524]" />
            {product.offers_count} {product.offers_count === 1 ? 'offer' : 'offers'}
          </span>
        </div>
      </div>

      {/* Card Body */}
      <div className="p-3.5 flex-1 flex flex-col justify-between gap-2.5">
        <div>
          {/* Region & Rating Row */}
          <div className="flex items-center justify-between gap-2 mb-1.5">
            <span
              className={`text-[10px] font-extrabold uppercase tracking-wider px-2 py-0.5 rounded border ${regionColorClass}`}
            >
              {product.region}
            </span>
            <span className="inline-flex items-center gap-1 text-xs font-bold text-[#ffd524]">
              <Star className="w-3.5 h-3.5 fill-[#ffd524]" />
              {product.rating.toFixed(2)}
            </span>
          </div>

          {/* Product Title */}
          <Link
            to={`/product/${product.slug}`}
            className="text-sm font-extrabold text-white group-hover:text-[#ffd524] line-clamp-2 leading-snug transition-colors"
          >
            {product.title}
          </Link>

          {product.matched_correction && (
            <div className="mt-1 text-[10px] font-bold text-[#ffd524] bg-[#ffd524]/10 px-2 py-0.5 rounded inline-block">
              ✓ Matched typo query
            </div>
          )}
        </div>

        {/* Pricing & Quick Add to Cart Row */}
        <div className="pt-2 border-t border-[#3b156b] flex items-end justify-between gap-2">
          <div>
            <span className="block text-[10px] uppercase tracking-wider text-purple-300 font-semibold">
              From {product.base_price > product.lowest_price && (
                <del className="text-purple-400/80 ml-1">{formatPrice(product.base_price)}</del>
              )}
            </span>
            <span className="text-lg font-extrabold text-white leading-none">
              {formatPrice(product.lowest_price)}
            </span>
          </div>

          <button
            type="button"
            onClick={() => addToCart(product)}
            title="Quick Add Best Offer to Cart"
            className="px-3 py-2 rounded-xl bg-[#ffd524] hover:bg-[#ffe055] text-[#16072c] font-extrabold text-xs flex items-center gap-1.5 shadow-md shadow-[#ffd524]/15 transition-transform active:scale-95 cursor-pointer"
          >
            <ShoppingCart className="w-3.5 h-3.5" />
            <span>Buy</span>
          </button>
        </div>
      </div>
    </div>
  );
};
