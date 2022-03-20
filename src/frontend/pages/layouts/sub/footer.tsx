import React from 'react';
import { Layout } from 'antd';

export function BaseFooter(props: any) {
  const { Footer } = Layout;

  return (
    <div>
      <Footer style={{ textAlign: 'center' }}>©2020 Created by GrowerLab</Footer>
    </div>
  );
}
