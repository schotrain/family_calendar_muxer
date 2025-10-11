import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest';
import { apiClient } from './apiClient';
import { authService } from './auth';

describe('ApiClient', () => {
  const mockFetch = vi.fn();
  const originalFetch = global.fetch;
  const originalLocation = window.location;

  beforeEach(() => {
    global.fetch = mockFetch;
    mockFetch.mockClear();
    sessionStorage.clear();

    // Mock window.location
    delete (window as any).location;
    window.location = { ...originalLocation, href: '' } as any;
  });

  afterEach(() => {
    global.fetch = originalFetch;
    window.location = originalLocation;
  });

  describe('get', () => {
    it('should make GET request with correct headers', async () => {
      const mockData = { id: 1, name: 'test' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => mockData,
      });

      const result = await apiClient.get('/api/test');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/test',
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
          credentials: 'include',
        })
      );
      expect(result).toEqual(mockData);
    });

    it('should include Authorization header when token exists', async () => {
      authService.setToken('test-token');
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => ({}),
      });

      await apiClient.get('/api/test');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/test',
        expect.objectContaining({
          headers: expect.objectContaining({
            'Authorization': 'Bearer test-token',
          }),
        })
      );
    });

    it('should throw error on non-ok response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        statusText: 'Internal Server Error',
      });

      await expect(apiClient.get('/api/test')).rejects.toThrow('API error: 500 Internal Server Error');
    });

    it('should redirect to home and remove token on 401 response', async () => {
      authService.setToken('invalid-token');
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
      });

      await expect(apiClient.get('/api/test')).rejects.toThrow('API error: 401 Unauthorized');
      expect(authService.getToken()).toBeNull();
      expect(window.location.href).toBe('/');
    });
  });

  describe('post', () => {
    it('should make POST request with data', async () => {
      const requestData = { name: 'test', value: 123 };
      const responseData = { id: 1, ...requestData };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      });

      const result = await apiClient.post('/api/test', requestData);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/test',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(requestData),
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
          }),
        })
      );
      expect(result).toEqual(responseData);
    });
  });

  describe('put', () => {
    it('should make PUT request with data', async () => {
      const requestData = { name: 'updated' };
      const responseData = { id: 1, ...requestData };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      });

      const result = await apiClient.put('/api/test/1', requestData);

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/test/1',
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(requestData),
        })
      );
      expect(result).toEqual(responseData);
    });
  });

  describe('delete', () => {
    it('should make DELETE request', async () => {
      const responseData = { message: 'deleted' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: async () => responseData,
      });

      const result = await apiClient.delete('/api/test/1');

      expect(mockFetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/test/1',
        expect.objectContaining({
          method: 'DELETE',
        })
      );
      expect(result).toEqual(responseData);
    });
  });
});
