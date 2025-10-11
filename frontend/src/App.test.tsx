import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import App from './App';

// Mock the pages since we just want to test routing
vi.mock('./pages/Home', () => ({
  default: () => <div>Home Page</div>,
}));

vi.mock('./pages/AuthCallback', () => ({
  default: () => <div>Auth Callback Page</div>,
}));

describe('App', () => {
  it('should render without crashing', () => {
    render(<App />);
    expect(screen.getByText('Home Page')).toBeInTheDocument();
  });
});
