import { apiClient } from '../utils/apiClient';
import type { User, HealthCheckResponse } from './types';

// User API
export const userApi = {
  getUserInfo: async (): Promise<User> => {
    return apiClient.get<User>('/api/userinfo');
  },
};

// Health API
export const healthApi = {
  check: async (): Promise<HealthCheckResponse> => {
    return apiClient.get<HealthCheckResponse>('/health');
  },
};

// Re-export types for convenience
export type { User, HealthCheckResponse };
