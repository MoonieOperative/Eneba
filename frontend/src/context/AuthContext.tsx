import React, { createContext, useContext, useEffect, useState } from 'react';
import type { User } from '../types/marketplace';
import { getStoredToken, marketplaceApi, setStoredToken } from '../services/api';

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<User>;
  register: (email: string, username: string, password: string) => Promise<User>;
  demoLogin: () => Promise<User>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
  setUser: React.Dispatch<React.SetStateAction<User | null>>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  const refreshUser = async () => {
    try {
      const res = await marketplaceApi.getMe();
      setUser(res.user);
    } catch {
      setUser(null);
      setStoredToken(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (getStoredToken()) {
      void refreshUser();
    } else {
      // Still attempt cookie-based session check
      void refreshUser();
    }
  }, []);

  const login = async (email: string, password: string) => {
    const res = await marketplaceApi.login(email, password);
    setUser(res.user);
    return res.user;
  };

  const register = async (email: string, username: string, password: string) => {
    const res = await marketplaceApi.register(email, username, password);
    setUser(res.user);
    return res.user;
  };

  const demoLogin = async () => {
    const res = await marketplaceApi.demoLogin();
    setUser(res.user);
    return res.user;
  };

  const logout = async () => {
    await marketplaceApi.logout();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        register,
        demoLogin,
        logout,
        refreshUser,
        setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return ctx;
}
