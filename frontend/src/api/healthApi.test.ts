import { describe, it, expect, vi, beforeEach } from 'vitest';
import { healthApi } from './healthApi';
import { apiClient } from '../utils/apiClient';

vi.mock('../utils/apiClient', () => ({
  apiClient: {
    get: vi.fn(),
  },
}));

describe('healthApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('check', () => {
    it('should fetch health status from /health', async () => {
      const mockHealth = {
        status: 'ok',
        message: 'Service is healthy',
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(mockHealth);

      const result = await healthApi.check();

      expect(apiClient.get).toHaveBeenCalledWith('/health');
      expect(result).toEqual(mockHealth);
    });
  });
});
