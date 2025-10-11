import { describe, it, expect, vi, beforeEach } from 'vitest';
import { userApi } from './userApi';
import { apiClient } from '../utils/apiClient';

vi.mock('../utils/apiClient', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

describe('userApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getUserInfo', () => {
    it('should fetch user info from /api/userinfo', async () => {
      const mockUser = {
        id: 1,
        given_name: 'John',
        family_name: 'Doe',
        email: 'john.doe@example.com',
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(mockUser);

      const result = await userApi.getUserInfo();

      expect(apiClient.get).toHaveBeenCalledWith('/api/userinfo');
      expect(result).toEqual(mockUser);
    });
  });
});
