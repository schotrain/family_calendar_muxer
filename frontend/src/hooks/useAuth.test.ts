import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useAuth } from './useAuth';
import { authService } from '../utils/auth';
import { userApi } from '../api';

vi.mock('../utils/auth', () => ({
  authService: {
    isAuthenticated: vi.fn(),
    removeToken: vi.fn(),
  },
}));

vi.mock('../api', () => ({
  userApi: {
    getUserInfo: vi.fn(),
  },
}));

describe('useAuth', () => {
  const mockUser = {
    id: 1,
    given_name: 'John',
    family_name: 'Doe',
    email: 'john@example.com',
  };

  beforeEach(() => {
    vi.clearAllMocks();
    delete (window as any).location;
    window.location = { href: '' } as any;
  });

  it('should return null user when not authenticated', async () => {
    vi.mocked(authService.isAuthenticated).mockReturnValue(false);

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.user).toBeNull();
    expect(result.current.isAuthenticated).toBe(false);
    expect(result.current.error).toBeNull();
  });

  it('should fetch and return user when authenticated', async () => {
    vi.mocked(authService.isAuthenticated).mockReturnValue(true);
    vi.mocked(userApi.getUserInfo).mockResolvedValue(mockUser);

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.user).toEqual(mockUser);
    expect(result.current.isAuthenticated).toBe(true);
    expect(result.current.error).toBeNull();
  });

  it('should handle error when fetching user fails', async () => {
    vi.mocked(authService.isAuthenticated).mockReturnValue(true);
    vi.mocked(userApi.getUserInfo).mockRejectedValue(new Error('Network error'));

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.user).toBeNull();
    expect(result.current.error).toBe('Failed to load user information');
    expect(authService.removeToken).toHaveBeenCalled();
  });

  it('should logout and redirect to home', async () => {
    vi.mocked(authService.isAuthenticated).mockReturnValue(true);
    vi.mocked(userApi.getUserInfo).mockResolvedValue(mockUser);

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    result.current.logout();

    expect(authService.removeToken).toHaveBeenCalled();
    expect(window.location.href).toBe('/');
  });
});
