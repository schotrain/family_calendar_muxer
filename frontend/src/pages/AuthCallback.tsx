import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Spin } from 'antd';
import { authService } from '../utils/auth';

const AuthCallback: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    const token = searchParams.get('token');

    if (token) {
      authService.setToken(token);
      navigate('/', { replace: true });
    } else {
      navigate('/', { replace: true });
    }
  }, [searchParams, navigate]);

  return (
    <div style={{
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      height: '100vh'
    }}>
      <Spin size="large" />
    </div>
  );
};

export default AuthCallback;
