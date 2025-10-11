import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import AppContent from './AppContent';

vi.mock('./CalendarMuxList', () => ({
  default: () => <div>Calendar Mux List</div>,
}));

describe('AppContent', () => {
  it('should render CalendarMuxList', () => {
    render(<AppContent />);
    expect(screen.getByText('Calendar Mux List')).toBeInTheDocument();
  });
});
