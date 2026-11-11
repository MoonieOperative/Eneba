import React from 'react';
import { Link } from 'react-router-dom';
import { ShieldCheck, Lock, Zap, Globe, Code2 } from 'lucide-react';

export const Footer: React.FC = () => {
  return (
    <footer className="bg-[#120422] border-t border-[#3b156b] mt-20 text-purple-200/80 text-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 pb-10 border-b border-[#2e1054]">
          {/* Brand Column */}
          <div className="space-y-3">
            <div className="flex items-center gap-2.5">
              <div className="w-8 h-8 rounded-lg bg-[#ffd524] flex items-center justify-center text-[#16072c] font-black">
                E
              </div>
              <span className="text-xl font-extrabold text-white tracking-tight">eneba</span>
            </div>
            <p className="text-xs text-purple-300 leading-relaxed">
              The fastest-growing digital marketplace for gamers. Instant game keys, DLCs, gift cards, and subscriptions with AES-256-GCM encrypted key vault delivery.
            </p>
            <div className="flex items-center gap-2 pt-1 text-[11px] text-[#00e676] font-bold">
              <ShieldCheck className="w-4 h-4" />
              <span>Eneba Shield Buyer Protection Active</span>
            </div>
          </div>

          {/* Store Links */}
          <div>
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-white mb-3">
              Marketplace Catalog
            </h3>
            <ul className="space-y-2 text-xs">
              <li><Link to="/store?platform=Steam" className="hover:text-[#ffd524]">Steam Game Keys</Link></li>
              <li><Link to="/store?platform=Xbox+Live" className="hover:text-[#ffd524]">Xbox Live & Game Pass</Link></li>
              <li><Link to="/store?platform=PlayStation" className="hover:text-[#ffd524]">PlayStation Network Cards</Link></li>
              <li><Link to="/store?platform=Nintendo+eShop" className="hover:text-[#ffd524]">Nintendo Switch eShop</Link></li>
              <li><Link to="/store?max_price=25" className="hover:text-[#ffd524]">Deals Under €25</Link></li>
            </ul>
          </div>

          {/* Account & Orders */}
          <div>
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-white mb-3">
              Buyer Account
            </h3>
            <ul className="space-y-2 text-xs">
              <li><Link to="/login" className="hover:text-[#ffd524]">Sign In / 1-Click Demo</Link></li>
              <li><Link to="/register" className="hover:text-[#ffd524]">Create Free Account</Link></li>
              <li><Link to="/orders" className="hover:text-[#ffd524]">My Orders & Digital Keys</Link></li>
              <li><Link to="/checkout" className="hover:text-[#ffd524]">Shopping Cart & Checkout</Link></li>
            </ul>
          </div>

          {/* Security & Architecture Info */}
          <div className="space-y-2.5">
            <h3 className="text-xs font-extrabold uppercase tracking-wider text-white mb-3">
              Security & Architecture
            </h3>
            <div className="p-3 rounded-xl bg-[#1d0837] border border-[#3e176e] text-xs space-y-1.5">
              <div className="flex items-center gap-1.5 text-[#ffd524] font-bold">
                <Lock className="w-3.5 h-3.5" />
                <span>Argon2id + AES-256-GCM</span>
              </div>
              <p className="text-[11px] text-purple-300">
                Zero plaintext credentials or keys stored at rest. All homepage blocks served dynamically by Go backend services.
              </p>
            </div>
          </div>
        </div>

        <div className="pt-6 flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-purple-400">
          <div className="flex items-center gap-4">
            <span className="flex items-center gap-1">
              <Code2 className="w-4 h-4 text-[#ffd524]" />
              Eneba Intern Assignment 2027 — React 19 + TypeScript + Go REST API + SQLite
            </span>
          </div>
          <div className="flex items-center gap-3 text-[11px] font-semibold">
            <span className="px-2.5 py-1 rounded bg-[#240a42] text-purple-200">VISA / MC Sandbox</span>
            <span className="px-2.5 py-1 rounded bg-[#240a42] text-purple-200">PayPal Test Mode</span>
            <span className="px-2.5 py-1 rounded bg-[#240a42] text-[#ffd524]">Eneba Purple Wallet</span>
            <span className="flex items-center gap-1 text-[#00e676]">
              <Zap className="w-3 h-3" /> Instant Key Dispatch
            </span>
            <span className="flex items-center gap-1">
              <Globe className="w-3 h-3" /> Global & EU
            </span>
          </div>
        </div>
      </div>
    </footer>
  );
};
