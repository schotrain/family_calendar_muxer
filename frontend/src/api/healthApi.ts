import { apiClient } from '../utils/apiClient';

export interface HealthCheckResponse {
  status: string;
  message: string;
}

export const healthApi = {
  check: async (): Promise<HealthCheckResponse> => {
    return apiClient.get<HealthCheckResponse>('/health');
  },
};
