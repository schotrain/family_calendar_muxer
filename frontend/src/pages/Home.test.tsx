import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import Home from './Home';

// Mock components
vi.mock('../components/AppHeader', () => ({
  default: () => <div>App Header</div>,
}));

vi.mock('../components/AppContent', () => ({
  default: () => <div>App Content</div>,
}));

describe('Home', () => {
  it('should render AppHeader and AppContent', () => {
    render(<Home />);
    expect(screen.getByText('App Header')).toBeInTheDocument();
    expect(screen.getByText('App Content')).toBeInTheDocument();
  });
});
