import { apiClient } from '../utils/apiClient';

export interface CalendarMux {
  id: number;
  created_by_id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface CalendarMuxListResponse {
  calendar_muxes: CalendarMux[];
}

export interface CreateCalendarMuxRequest {
  name: string;
  description: string;
}

export interface DeleteCalendarMuxResponse {
  message: string;
}

export const calendarMuxApi = {
  list: async (): Promise<CalendarMuxListResponse> => {
    return apiClient.get<CalendarMuxListResponse>('/api/calendar-mux');
  },

  create: async (data: CreateCalendarMuxRequest): Promise<CalendarMux> => {
    return apiClient.post<CalendarMux>('/api/calendar-mux', data);
  },

  delete: async (id: number): Promise<DeleteCalendarMuxResponse> => {
    return apiClient.delete<DeleteCalendarMuxResponse>(`/api/calendar-mux/${id}`);
  },
};
