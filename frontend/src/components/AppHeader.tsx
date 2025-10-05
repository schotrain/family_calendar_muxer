import React from 'react';
import { Layout, Button, Space, Typography } from 'antd';
import { useAuth } from '../hooks/useAuth';

const { Header } = Layout;
const { Text } = Typography;

const AppHeader: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();

  const handleSignIn = () => {
    const authUrl = import.meta.env.PUBLIC_AUTH_LOGIN_URL || 'http://localhost:8080/auth/google';
    const callbackUrl = `${window.location.origin}/auth/callback`;
    window.location.href = `${authUrl}?callback=${encodeURIComponent(callbackUrl)}`;
  };

  return (
    <Header
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        background: '#001529',
        padding: '0 50px',
      }}
    >
      <div
        style={{
          color: 'white',
          fontSize: '20px',
          fontWeight: 'bold',
        }}
      >
        Family Calendar Muxer
      </div>
      {isAuthenticated && user ? (
        <Space>
          <Text style={{ color: 'white' }}>Hi {user.given_name}</Text>
          <Button onClick={logout}>Sign Out</Button>
        </Space>
      ) : (
        <Button type="primary" onClick={handleSignIn}>
          Sign In
        </Button>
      )}
    </Header>
  );
};

export default AppHeader;
