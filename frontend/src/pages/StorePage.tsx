import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import {
  SlidersHorizontal,
  Search,
  RotateCcw,
  Sparkles,
  ChevronLeft,
  ChevronRight,
  Check,
  X,
} from 'lucide-react';
import { marketplaceApi } from '../services/api';
import type { CatalogResponse } from '../types/marketplace';
import { ProductCard } from '../components/ProductCard';

export const StorePage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();

  // Read filters from URL search params
  const queryParam = searchParams.get('q') || '';
  const platformParam = searchParams.get('platform') ? searchParams.get('platform')!.split(',') : [];
  const regionParam = searchParams.get('region') ? searchParams.get('region')!.split(',') : [];
  const genreParam = searchParams.get('genre') ? searchParams.get('genre')!.split(',') : [];
  const productTypeParam = searchParams.get('product_type') || '';
  const minPriceParam = searchParams.get('min_price') || '';
  const maxPriceParam = searchParams.get('max_price') || '';
  const inStockParam = searchParams.get('in_stock') === 'true';
  const sortParam = searchParams.get('sort') || '';
  const pageParam = Number(searchParams.get('page') || '1');

  const [localSearch, setLocalSearch] = useState(queryParam);
  const [minPriceInput, setMinPriceInput] = useState(minPriceParam);
  const [maxPriceInput, setMaxPriceInput] = useState(maxPriceParam);
  const [catalog, setCatalog] = useState<CatalogResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false);

  useEffect(() => {
    setLocalSearch(queryParam);
  }, [queryParam]);

  useEffect(() => {
    setMinPriceInput(minPriceParam);
    setMaxPriceInput(maxPriceParam);
  }, [minPriceParam, maxPriceParam]);

  useEffect(() => {
    let active = true;
    setLoading(true);
    marketplaceApi
      .getCatalog({
        q: queryParam,
        platform: platformParam,
        region: regionParam,
        genre: genreParam,
        product_type: productTypeParam,
        min_price: minPriceParam,
        max_price: maxPriceParam,
        in_stock: inStockParam,
        sort: sortParam,
        page: pageParam,
        page_size: 12,
      })
      .then((res) => {
        if (active) setCatalog(res);
      })
      .finally(() => {
        if (active) setLoading(false);
      });

    return () => {
      active = false;
    };
  }, [searchParams]);

  const updateParam = (key: string, value: string | null) => {
    const next = new URLSearchParams(searchParams);
    if (value && value.trim() !== '') {
      next.set(key, value);
    } else {
      next.delete(key);
    }
    if (key !== 'page') {
      next.set('page', '1');
    }
    setSearchParams(next);
  };

  const toggleMultiParam = (key: string, currentList: string[], item: string) => {
    const exists = currentList.some((c) => c.toLowerCase() === item.toLowerCase());
    const updated = exists
      ? currentList.filter((c) => c.toLowerCase() !== item.toLowerCase())
      : [...currentList, item];
    updateParam(key, updated.length > 0 ? updated.join(',') : null);
  };

  const handleApplyPrice = (e: React.FormEvent) => {
    e.preventDefault();
    const next = new URLSearchParams(searchParams);
    if (minPriceInput.trim()) next.set('min_price', minPriceInput.trim());
    else next.delete('min_price');
    if (maxPriceInput.trim()) next.set('max_price', maxPriceInput.trim());
    else next.delete('max_price');
    next.set('page', '1');
    setSearchParams(next);
  };

  const resetAllFilters = () => {
    setLocalSearch('');
    setMinPriceInput('');
    setMaxPriceInput('');
    setSearchParams(new URLSearchParams());
  };

  const hasActiveFilters =
    queryParam ||
    platformParam.length > 0 ||
    regionParam.length > 0 ||
    genreParam.length > 0 ||
    productTypeParam ||
    minPriceParam ||
    maxPriceParam ||
    inStockParam;

  const platforms = catalog?.available_facets?.platforms || [];
  const regions = catalog?.available_facets?.regions || [];
  const genres = catalog?.available_facets?.genres || [];
  const productTypes = catalog?.available_facets?.product_types || [];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8">
      {/* Top Store Header */}
      <div className="mb-6 flex flex-col md:flex-row md:items-end justify-between gap-4 pb-5 border-b border-[#3b156b]">
        <div>
          <span className="text-xs font-extrabold uppercase tracking-wider text-[#ffd524]">
            Eneba Official Marketplace
          </span>
          <h1 className="text-2xl sm:text-3xl font-black text-white mt-0.5">
            {queryParam ? `Search Results for "${queryParam}"` : 'Browse All Digital Games & Keys'}
          </h1>
          <p className="text-xs text-purple-300 mt-1">
            Showing {catalog?.total_count ?? 0} verified instant-delivery products across global sellers
          </p>
        </div>

        {/* Sort & Mobile Filter Controls */}
        <div className="flex items-center gap-3 flex-wrap">
          <button
            type="button"
            onClick={() => setMobileFiltersOpen(!mobileFiltersOpen)}
            className="lg:hidden inline-flex items-center gap-2 px-4 py-2.5 rounded-xl bg-[#311159] border border-[#542294] text-xs font-bold text-white cursor-pointer"
          >
            <SlidersHorizontal className="w-4 h-4 text-[#ffd524]" />
            <span>Filters</span>
          </button>

          <div className="flex items-center gap-2 bg-[#230942] border border-[#4c1e86] rounded-xl px-3.5 py-2">
            <span className="text-xs text-purple-300 font-semibold">Sort by:</span>
            <select
              aria-label="Sort products"
              value={sortParam}
              onChange={(e) => updateParam('sort', e.target.value)}
              className="bg-transparent text-xs font-extrabold text-[#ffd524] focus:outline-none cursor-pointer"
            >
              <option value="" className="bg-[#1f083d] text-white">
                {queryParam ? 'Best Search Match' : 'Most Popular'}
              </option>
              <option value="popular" className="bg-[#1f083d] text-white">Best Selling</option>
              <option value="price_asc" className="bg-[#1f083d] text-white">Price: Low to High</option>
              <option value="price_desc" className="bg-[#1f083d] text-white">Price: High to Low</option>
              <option value="discount_desc" className="bg-[#1f083d] text-white">Biggest Discount (%)</option>
              <option value="rating_desc" className="bg-[#1f083d] text-white">Highest User Rating</option>
              <option value="release_desc" className="bg-[#1f083d] text-white">Newest Release Date</option>
            </select>
          </div>
        </div>
      </div>

      {/* Typo-Tolerance Auto-Match Callout */}
      {catalog?.suggested_query && queryParam && (
        <div className="mb-6 p-4 rounded-2xl bg-[#ffd524]/15 border border-[#ffd524]/40 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div className="flex items-center gap-2.5 text-sm text-purple-100">
            <Sparkles className="w-5 h-5 text-[#ffd524] shrink-0" />
            <span>
              Our typo-tolerant search matched <strong className="text-white">"{queryParam}"</strong> to{' '}
              <strong className="text-[#ffd524]">{catalog.suggested_query}</strong>
            </span>
          </div>
          <button
            type="button"
            onClick={() => updateParam('q', catalog.suggested_query!)}
            className="text-xs font-extrabold px-3 py-1.5 rounded-lg bg-[#ffd524] text-[#16072c] self-start sm:self-auto cursor-pointer"
          >
            Search Exact Title
          </button>
        </div>
      )}

      {/* Active Filter Chips */}
      {hasActiveFilters && (
        <div className="mb-6 flex items-center gap-2 flex-wrap">
          <span className="text-xs font-bold text-purple-300 mr-1">Active Filters:</span>
          {queryParam && (
            <FilterChip label={`Search: "${queryParam}"`} onRemove={() => updateParam('q', null)} />
          )}
          {platformParam.map((p) => (
            <FilterChip
              key={p}
              label={`Platform: ${p}`}
              onRemove={() => toggleMultiParam('platform', platformParam, p)}
            />
          ))}
          {regionParam.map((r) => (
            <FilterChip
              key={r}
              label={`Region: ${r}`}
              onRemove={() => toggleMultiParam('region', regionParam, r)}
            />
          ))}
          {genreParam.map((g) => (
            <FilterChip
              key={g}
              label={`Genre: ${g}`}
              onRemove={() => toggleMultiParam('genre', genreParam, g)}
            />
          ))}
          {productTypeParam && (
            <FilterChip
              label={`Type: ${productTypeParam}`}
              onRemove={() => updateParam('product_type', null)}
            />
          )}
          {(minPriceParam || maxPriceParam) && (
            <FilterChip
              label={`Price: €${minPriceParam || '0'} - €${maxPriceParam || 'Any'}`}
              onRemove={() => {
                const next = new URLSearchParams(searchParams);
                next.delete('min_price');
                next.delete('max_price');
                setSearchParams(next);
              }}
            />
          )}
          <button
            type="button"
            onClick={resetAllFilters}
            className="inline-flex items-center gap-1 text-xs font-extrabold text-[#ff4081] hover:underline ml-2 cursor-pointer"
          >
            <RotateCcw className="w-3.5 h-3.5" />
            Reset All
          </button>
        </div>
      )}

      {/* Main Layout: Sidebar Filters + Product Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left Filter Sidebar */}
        <aside
          className={`lg:col-span-3 space-y-6 bg-[#220940] border border-[#481d7d] rounded-2xl p-5 ${
            mobileFiltersOpen ? 'block' : 'hidden lg:block'
          }`}
        >
          {/* In-Sidebar Keyword Filter */}
          <div className="space-y-2">
            <label className="text-xs font-extrabold uppercase tracking-wider text-purple-200 block">
              Filter by Name (Typo-Tolerant)
            </label>
            <form
              onSubmit={(e) => {
                e.preventDefault();
                updateParam('q', localSearch);
              }}
              className="relative"
            >
              <input
                type="text"
                value={localSearch}
                onChange={(e) => setLocalSearch(e.target.value)}
                placeholder="e.g. witcer 3, eldn ring..."
                className="w-full pl-3 pr-9 py-2 rounded-xl bg-[#16062b] border border-[#4d1f85] focus:border-[#ffd524] text-xs text-white focus:outline-none"
              />
              <button
                type="submit"
                className="absolute right-2 top-1/2 -translate-y-1/2 text-[#ffd524] cursor-pointer"
                aria-label="Apply search filter"
              >
                <Search className="w-4 h-4" />
              </button>
            </form>
          </div>

          {/* Price Range Filter */}
          <div className="pt-4 border-t border-[#3b156b] space-y-3">
            <span className="text-xs font-extrabold uppercase tracking-wider text-purple-200 block">
              Price Range (EUR)
            </span>
            <form onSubmit={handleApplyPrice} className="space-y-2.5">
              <div className="grid grid-cols-2 gap-2">
                <input
                  type="number"
                  min="0"
                  step="0.5"
                  placeholder="Min €"
                  value={minPriceInput}
                  onChange={(e) => setMinPriceInput(e.target.value)}
                  className="w-full px-3 py-2 rounded-xl bg-[#16062b] border border-[#4d1f85] text-xs text-white focus:border-[#ffd524] focus:outline-none"
                />
                <input
                  type="number"
                  min="0"
                  step="0.5"
                  placeholder="Max €"
                  value={maxPriceInput}
                  onChange={(e) => setMaxPriceInput(e.target.value)}
                  className="w-full px-3 py-2 rounded-xl bg-[#16062b] border border-[#4d1f85] text-xs text-white focus:border-[#ffd524] focus:outline-none"
                />
              </div>
              <button
                type="submit"
                className="w-full py-2 rounded-xl bg-[#381366] hover:bg-[#ffd524] text-white hover:text-[#16072c] font-extrabold text-xs transition-colors cursor-pointer"
              >
                Apply Price Range
              </button>
            </form>

            {/* Quick Budget Presets */}
            <div className="flex flex-wrap gap-1.5 pt-1">
              {[15, 25, 40].map((maxVal) => (
                <button
                  key={maxVal}
                  type="button"
                  onClick={() => updateParam('max_price', String(maxVal))}
                  className={`text-[11px] font-bold px-2.5 py-1 rounded-lg border transition-colors cursor-pointer ${
                    maxPriceParam === String(maxVal)
                      ? 'bg-[#ffd524] text-[#16072c] border-[#ffd524]'
                      : 'bg-[#16062b] text-purple-200 border-[#441b77] hover:border-[#ffd524]'
                  }`}
                >
                  Under €{maxVal}
                </button>
              ))}
            </div>
          </div>

          {/* Platform Filter */}
          <div className="pt-4 border-t border-[#3b156b] space-y-2.5">
            <span className="text-xs font-extrabold uppercase tracking-wider text-purple-200 block">
              Platform / DRM
            </span>
            <div className="space-y-1.5">
              {platforms.map((f) => {
                const checked = platformParam.some((p) => p.toLowerCase() === f.value.toLowerCase());
                return (
                  <FacetCheckbox
                    key={f.value}
                    label={f.value}
                    count={f.count}
                    checked={checked}
                    onChange={() => toggleMultiParam('platform', platformParam, f.value)}
                  />
                );
              })}
            </div>
          </div>

          {/* Activation Region Filter */}
          <div className="pt-4 border-t border-[#3b156b] space-y-2.5">
            <span className="text-xs font-extrabold uppercase tracking-wider text-purple-200 block">
              Activation Region
            </span>
            <div className="space-y-1.5">
              {regions.map((f) => {
                const checked = regionParam.some((r) => r.toLowerCase() === f.value.toLowerCase());
                return (
                  <FacetCheckbox
                    key={f.value}
                    label={f.value}
                    count={f.count}
                    checked={checked}
                    onChange={() => toggleMultiParam('region', regionParam, f.value)}
                  />
                );
              })}
            </div>
          </div>

          {/* Product Type Filter */}
          <div className="pt-4 border-t border-[#3b156b] space-y-2.5">
            <span className="text-xs font-extrabold uppercase tracking-wider text-purple-200 block">
              Product Type
            </span>
            <div className="space-y-1.5">
              {productTypes.map((f) => {
                const checked = productTypeParam.toLowerCase() === f.value.toLowerCase();
                return (
                  <FacetCheckbox
                    key={f.value}
                    label={f.value}
                    count={f.count}
                    checked={checked}
                    onChange={() => updateParam('product_type', checked ? null : f.value)}
                  />
                );
              })}
            </div>
          </div>

          {/* Genre Filter */}
          <div className="pt-4 border-t border-[#3b156b] space-y-2.5">
            <span className="text-xs font-extrabold uppercase tracking-wider text-purple-200 block">
              Genre
            </span>
            <div className="space-y-1.5">
              {genres.map((f) => {
                const checked = genreParam.some((g) => g.toLowerCase() === f.value.toLowerCase());
                return (
                  <FacetCheckbox
                    key={f.value}
                    label={f.value}
                    count={f.count}
                    checked={checked}
                    onChange={() => toggleMultiParam('genre', genreParam, f.value)}
                  />
                );
              })}
            </div>
          </div>
        </aside>

        {/* Right Product Catalog Grid */}
        <div className="lg:col-span-9 space-y-8">
          {loading ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-5">
              {Array.from({ length: 8 }).map((_, idx) => (
                <div
                  key={idx}
                  className="aspect-[3/4] rounded-2xl bg-[#260c47] animate-pulse border border-[#471c7c]"
                />
              ))}
            </div>
          ) : !catalog || catalog.products.length === 0 ? (
            <div className="rounded-3xl bg-[#240a42] border border-[#4c1e86] p-12 text-center space-y-4">
              <h3 className="text-xl font-extrabold text-white">No Matching Products Found</h3>
              <p className="text-sm text-purple-300 max-w-md mx-auto">
                We couldn't find any products matching all of your selected filters. Try broadening your price range or clearing active filters.
              </p>
              <button
                type="button"
                onClick={resetAllFilters}
                className="px-6 py-3 rounded-xl bg-[#ffd524] text-[#16072c] font-extrabold text-xs uppercase tracking-wider cursor-pointer"
              >
                Reset All Filters
              </button>
            </div>
          ) : (
            <>
              <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-5">
                {catalog.products.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>

              {/* Pagination Bar */}
              {catalog.total_pages > 1 && (
                <div className="pt-6 border-t border-[#3b156b] flex items-center justify-center gap-2">
                  <button
                    type="button"
                    disabled={catalog.page <= 1}
                    onClick={() => updateParam('page', String(catalog.page - 1))}
                    className="p-2.5 rounded-xl bg-[#270c47] border border-[#491d7d] text-white disabled:opacity-40 hover:border-[#ffd524] cursor-pointer"
                  >
                    <ChevronLeft className="w-4 h-4" />
                  </button>
                  {Array.from({ length: catalog.total_pages }).map((_, idx) => {
                    const pNum = idx + 1;
                    return (
                      <button
                        key={pNum}
                        type="button"
                        onClick={() => updateParam('page', String(pNum))}
                        className={`w-10 h-10 rounded-xl font-extrabold text-xs transition-colors cursor-pointer ${
                          catalog.page === pNum
                            ? 'bg-[#ffd524] text-[#16072c]'
                            : 'bg-[#270c47] text-white border border-[#491d7d] hover:border-[#ffd524]'
                        }`}
                      >
                        {pNum}
                      </button>
                    );
                  })}
                  <button
                    type="button"
                    disabled={catalog.page >= catalog.total_pages}
                    onClick={() => updateParam('page', String(catalog.page + 1))}
                    className="p-2.5 rounded-xl bg-[#270c47] border border-[#491d7d] text-white disabled:opacity-40 hover:border-[#ffd524] cursor-pointer"
                  >
                    <ChevronRight className="w-4 h-4" />
                  </button>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

const FacetCheckbox: React.FC<{
  label: string;
  count: number;
  checked: boolean;
  onChange: () => void;
}> = ({ label, count, checked, onChange }) => (
  <button
    type="button"
    onClick={onChange}
    className="w-full flex items-center justify-between gap-2 py-1.5 px-2 rounded-lg hover:bg-[#2f0f57] transition-colors text-left cursor-pointer group"
  >
    <div className="flex items-center gap-2.5 min-w-0">
      <div
        className={`w-4 h-4 rounded flex items-center justify-center border transition-colors shrink-0 ${
          checked
            ? 'bg-[#ffd524] border-[#ffd524] text-[#16072c]'
            : 'bg-[#16062b] border-[#592499] group-hover:border-[#ffd524]'
        }`}
      >
        {checked && <Check className="w-3 h-3 stroke-[3]" />}
      </div>
      <span className={`text-xs truncate ${checked ? 'font-bold text-white' : 'text-purple-200'}`}>
        {label}
      </span>
    </div>
    <span className="text-[11px] font-mono text-purple-400 shrink-0">{count}</span>
  </button>
);

const FilterChip: React.FC<{ label: string; onRemove: () => void }> = ({ label, onRemove }) => (
  <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#32115c] border border-[#5b259e] text-xs font-bold text-white">
    <span>{label}</span>
    <button
      type="button"
      onClick={onRemove}
      className="text-purple-300 hover:text-[#ffd524] cursor-pointer"
      aria-label={`Remove ${label}`}
    >
      <X className="w-3.5 h-3.5" />
    </button>
  </span>
);
