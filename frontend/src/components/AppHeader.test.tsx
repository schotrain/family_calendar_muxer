import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import AppHeader from './AppHeader';
import { useAuth } from '../hooks/useAuth';

vi.mock('../hooks/useAuth');

describe('AppHeader', () => {
  const mockUser = {
    id: 1,
    given_name: 'John',
    family_name: 'Doe',
    email: 'john@example.com',
  };

  beforeEach(() => {
    vi.clearAllMocks();
    delete (window as any).location;
    window.location = { href: '', origin: 'http://localhost:3000' } as any;
  });

  it('should render sign in button when not authenticated', () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      loading: false,
      error: null,
      isAuthenticated: false,
      logout: vi.fn(),
    });

    render(<AppHeader />);
    expect(screen.getByText('Sign In')).toBeInTheDocument();
  });

  it('should render user name and sign out button when authenticated', () => {
    const mockLogout = vi.fn();
    vi.mocked(useAuth).mockReturnValue({
      user: mockUser,
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: mockLogout,
    });

    render(<AppHeader />);
    expect(screen.getByText('Hi John')).toBeInTheDocument();
    expect(screen.getByText('Sign Out')).toBeInTheDocument();
  });

  it('should call logout when sign out button is clicked', () => {
    const mockLogout = vi.fn();
    vi.mocked(useAuth).mockReturnValue({
      user: mockUser,
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: mockLogout,
    });

    render(<AppHeader />);
    fireEvent.click(screen.getByText('Sign Out'));
    expect(mockLogout).toHaveBeenCalled();
  });

  it('should redirect to auth URL when sign in button is clicked', () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      loading: false,
      error: null,
      isAuthenticated: false,
      logout: vi.fn(),
    });

    render(<AppHeader />);
    fireEvent.click(screen.getByText('Sign In'));

    expect(window.location.href).toContain('/auth/google');
    expect(window.location.href).toContain('callback=');
  });
});
