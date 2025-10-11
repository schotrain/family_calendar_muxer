import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import CalendarMuxList from './CalendarMuxList';
import { useAuth } from '../hooks/useAuth';
import { calendarMuxApi } from '../api';
import { message } from 'antd';

vi.mock('../hooks/useAuth');
vi.mock('../api');
vi.mock('antd', async () => {
  const actual = await vi.importActual('antd');
  return {
    ...actual,
    message: {
      error: vi.fn(),
      success: vi.fn(),
    },
  };
});

describe('CalendarMuxList', () => {
  const mockCalendarMuxes = [
    {
      id: 1,
      created_by_id: 1,
      name: 'Test Calendar 1',
      description: 'Description 1',
      created_at: '2025-01-01T00:00:00Z',
      updated_at: '2025-01-01T00:00:00Z',
    },
    {
      id: 2,
      created_by_id: 1,
      name: 'Test Calendar 2',
      description: 'Description 2',
      created_at: '2025-01-02T00:00:00Z',
      updated_at: '2025-01-02T00:00:00Z',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should show login message when not authenticated', () => {
    vi.mocked(useAuth).mockReturnValue({
      user: null,
      loading: false,
      error: null,
      isAuthenticated: false,
      logout: vi.fn(),
    });

    render(<CalendarMuxList />);
    expect(screen.getByText('Please log in to view your calendar muxes')).toBeInTheDocument();
  });

  it('should fetch and display calendar muxes when authenticated', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockResolvedValue({
      calendar_muxes: mockCalendarMuxes,
    });

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(screen.getByText('Test Calendar 1')).toBeInTheDocument();
      expect(screen.getByText('Test Calendar 2')).toBeInTheDocument();
    });
  });

  it('should show error message when fetching fails', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockRejectedValue(new Error('Network error'));

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(message.error).toHaveBeenCalledWith('Failed to load calendar muxes');
    });
  });

  it('should open create modal when create button is clicked', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockResolvedValue({
      calendar_muxes: [],
    });

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(screen.getByText('Create Calendar Mux')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Calendar Mux'));

    await waitFor(() => {
      expect(screen.getByText('Name')).toBeInTheDocument();
    });
  });

  it('should create calendar mux when form is submitted', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockResolvedValue({
      calendar_muxes: [],
    });

    const newCalendarMux = {
      id: 3,
      created_by_id: 1,
      name: 'New Calendar',
      description: 'New Description',
      created_at: '2025-01-03T00:00:00Z',
      updated_at: '2025-01-03T00:00:00Z',
    };

    vi.mocked(calendarMuxApi.create).mockResolvedValue(newCalendarMux);

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(screen.getByText('Create Calendar Mux')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Calendar Mux'));

    await waitFor(() => {
      expect(screen.getByPlaceholderText('Enter calendar mux name')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByPlaceholderText('Enter calendar mux name'), {
      target: { value: 'New Calendar' },
    });

    fireEvent.change(screen.getByPlaceholderText('Enter calendar mux description'), {
      target: { value: 'New Description' },
    });

    const okButton = screen.getAllByRole('button').find(btn => btn.textContent === 'OK');
    if (okButton) {
      fireEvent.click(okButton);
    }

    await waitFor(() => {
      expect(calendarMuxApi.create).toHaveBeenCalledWith({
        name: 'New Calendar',
        description: 'New Description',
      });
      expect(message.success).toHaveBeenCalledWith('Calendar mux created successfully');
    });
  });

  it('should delete calendar mux when delete button is clicked', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockResolvedValue({
      calendar_muxes: mockCalendarMuxes,
    });

    vi.mocked(calendarMuxApi.delete).mockResolvedValue({
      message: 'Calendar mux deleted successfully',
    });

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(screen.getByText('Test Calendar 1')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByText('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(calendarMuxApi.delete).toHaveBeenCalledWith(1);
      expect(message.success).toHaveBeenCalledWith('Calendar mux deleted successfully');
    });
  });

  it('should show error message when delete fails', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockResolvedValue({
      calendar_muxes: mockCalendarMuxes,
    });

    vi.mocked(calendarMuxApi.delete).mockRejectedValue(new Error('Delete failed'));

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(screen.getByText('Test Calendar 1')).toBeInTheDocument();
    });

    const deleteButtons = screen.getAllByText('Delete');
    fireEvent.click(deleteButtons[0]);

    await waitFor(() => {
      expect(message.error).toHaveBeenCalledWith('Failed to delete calendar mux');
    });
  });

  it('should show error message when create fails', async () => {
    vi.mocked(useAuth).mockReturnValue({
      user: { id: 1, given_name: 'John', family_name: 'Doe', email: 'john@example.com' },
      loading: false,
      error: null,
      isAuthenticated: true,
      logout: vi.fn(),
    });

    vi.mocked(calendarMuxApi.list).mockResolvedValue({
      calendar_muxes: [],
    });

    vi.mocked(calendarMuxApi.create).mockRejectedValue(new Error('Create failed'));

    render(<CalendarMuxList />);

    await waitFor(() => {
      expect(screen.getByText('Create Calendar Mux')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('Create Calendar Mux'));

    await waitFor(() => {
      expect(screen.getByPlaceholderText('Enter calendar mux name')).toBeInTheDocument();
    });

    fireEvent.change(screen.getByPlaceholderText('Enter calendar mux name'), {
      target: { value: 'New Calendar' },
    });

    const okButton = screen.getAllByRole('button').find(btn => btn.textContent === 'OK');
    if (okButton) {
      fireEvent.click(okButton);
    }

    await waitFor(() => {
      expect(message.error).toHaveBeenCalledWith('Failed to create calendar mux');
    });
  });
});
