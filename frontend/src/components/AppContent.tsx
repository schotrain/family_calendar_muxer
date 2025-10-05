import React from 'react';
import { Layout } from 'antd';

const { Content } = Layout;

const AppContent: React.FC = () => {
  return (
    <Content
      style={{
        padding: '50px',
        background: '#f0f2f5',
      }}
    >
      <div
        style={{
          background: 'white',
          padding: '24px',
          minHeight: '280px',
        }}
      >
        Content goes here
      </div>
    </Content>
  );
};

export default AppContent;
