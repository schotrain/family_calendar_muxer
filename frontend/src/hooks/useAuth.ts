import { useState, useEffect } from 'react';
import { authService } from '../utils/auth';
import { userApi, type User } from '../api';

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchUser = async () => {
      if (!authService.isAuthenticated()) {
        setLoading(false);
        return;
      }

      try {
        const userData = await userApi.getUserInfo();
        setUser(userData);
      } catch (err) {
        console.error('Failed to fetch user info:', err);
        setError('Failed to load user information');
        authService.removeToken();
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

  const logout = () => {
    authService.removeToken();
    setUser(null);
    window.location.href = '/';
  };

  return {
    user,
    loading,
    error,
    isAuthenticated: !!user,
    logout,
  };
};
