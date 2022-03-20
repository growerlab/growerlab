import React from 'react';
import { Layout } from 'antd';
import GQLProvider from '../../api/graphql/provider';
import { BaseHeader } from './sub/header';
import { BaseFooter } from './sub/footer';

export default function IndexLayout(props: any) {
  const { Content } = Layout;

  return (
    <GQLProvider>
      <Layout className="layout">
        <BaseHeader />
        <Content style={{ padding: '30px 50px' }}>
          <div className="site-layout-content">{props.children}</div>
        </Content>
        <BaseFooter />
      </Layout>
    </GQLProvider>
  );
}
