import { describe, it, expect, beforeEach, vi } from 'vitest';
import { authService } from './auth';

describe('authService', () => {
  beforeEach(() => {
    sessionStorage.clear();
  });

  describe('setToken', () => {
    it('should store token in sessionStorage', () => {
      const token = 'test-token-123';
      authService.setToken(token);
      expect(sessionStorage.getItem('auth_token')).toBe(token);
    });
  });

  describe('getToken', () => {
    it('should return token from sessionStorage', () => {
      const token = 'test-token-456';
      sessionStorage.setItem('auth_token', token);
      expect(authService.getToken()).toBe(token);
    });

    it('should return null when no token exists', () => {
      expect(authService.getToken()).toBeNull();
    });
  });

  describe('removeToken', () => {
    it('should remove token from sessionStorage', () => {
      sessionStorage.setItem('auth_token', 'test-token');
      authService.removeToken();
      expect(sessionStorage.getItem('auth_token')).toBeNull();
    });
  });

  describe('isAuthenticated', () => {
    it('should return true when token exists', () => {
      sessionStorage.setItem('auth_token', 'test-token');
      expect(authService.isAuthenticated()).toBe(true);
    });

    it('should return false when no token exists', () => {
      expect(authService.isAuthenticated()).toBe(false);
    });

    it('should return false when token is empty string', () => {
      sessionStorage.setItem('auth_token', '');
      expect(authService.isAuthenticated()).toBe(false);
    });
  });
});
