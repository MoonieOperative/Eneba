import { describe, expect, it } from 'vitest';
import { getStoredToken, setStoredToken } from './api';

describe('Marketplace Frontend API Service', () => {
  it('safely handles token storage operations', () => {
    setStoredToken(null);
    expect(getStoredToken()).toBeNull();
  });
});
