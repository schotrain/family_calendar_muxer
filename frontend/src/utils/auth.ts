const TOKEN_KEY = 'auth_token';

export const authService = {
  setToken(token: string): void {
    sessionStorage.setItem(TOKEN_KEY, token);
  },

  getToken(): string | null {
    return sessionStorage.getItem(TOKEN_KEY);
  },

  removeToken(): void {
    sessionStorage.removeItem(TOKEN_KEY);
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};
