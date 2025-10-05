import React from 'react';
import { Layout } from 'antd';
import AppHeader from '../components/AppHeader';
import AppContent from '../components/AppContent';

const Home: React.FC = () => {
  return (
    <Layout style={{ minHeight: '100vh' }}>
      <AppHeader />
      <AppContent />
    </Layout>
  );
};

export default Home;
