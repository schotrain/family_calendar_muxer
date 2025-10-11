import { apiClient } from '../utils/apiClient';

export interface User {
  id: number;
  given_name: string;
  family_name: string;
  email: string;
}

export const userApi = {
  getUserInfo: async (): Promise<User> => {
    return apiClient.get<User>('/api/userinfo');
  },
};
