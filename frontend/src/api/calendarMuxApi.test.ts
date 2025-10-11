import { describe, it, expect, vi, beforeEach } from 'vitest';
import { calendarMuxApi } from './calendarMuxApi';
import { apiClient } from '../utils/apiClient';

vi.mock('../utils/apiClient', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    delete: vi.fn(),
  },
}));

describe('calendarMuxApi', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('list', () => {
    it('should fetch calendar muxes from /api/calendar-mux', async () => {
      const mockResponse = {
        calendar_muxes: [
          {
            id: 1,
            created_by_id: 1,
            name: 'Test Calendar',
            description: 'Test Description',
            created_at: '2025-01-01T00:00:00Z',
            updated_at: '2025-01-01T00:00:00Z',
          },
        ],
      };

      vi.mocked(apiClient.get).mockResolvedValueOnce(mockResponse);

      const result = await calendarMuxApi.list();

      expect(apiClient.get).toHaveBeenCalledWith('/api/calendar-mux');
      expect(result).toEqual(mockResponse);
    });
  });

  describe('create', () => {
    it('should create calendar mux at /api/calendar-mux', async () => {
      const requestData = {
        name: 'New Calendar',
        description: 'New Description',
      };

      const mockResponse = {
        id: 2,
        created_by_id: 1,
        name: 'New Calendar',
        description: 'New Description',
        created_at: '2025-01-02T00:00:00Z',
        updated_at: '2025-01-02T00:00:00Z',
      };

      vi.mocked(apiClient.post).mockResolvedValueOnce(mockResponse);

      const result = await calendarMuxApi.create(requestData);

      expect(apiClient.post).toHaveBeenCalledWith('/api/calendar-mux', requestData);
      expect(result).toEqual(mockResponse);
    });
  });

  describe('delete', () => {
    it('should delete calendar mux at /api/calendar-mux/:id', async () => {
      const mockResponse = {
        message: 'Calendar mux deleted successfully',
      };

      vi.mocked(apiClient.delete).mockResolvedValueOnce(mockResponse);

      const result = await calendarMuxApi.delete(1);

      expect(apiClient.delete).toHaveBeenCalledWith('/api/calendar-mux/1');
      expect(result).toEqual(mockResponse);
    });
  });
});
