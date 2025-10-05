// API Types matching backend response schemas

export interface User {
  id: number;
  given_name: string;
  family_name: string;
  email: string;
}

export interface HealthCheckResponse {
  status: string;
  message: string;
}
