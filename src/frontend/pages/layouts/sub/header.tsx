import React from 'react';
import { Layout, Menu } from 'antd';
import Link from 'umi/link';
import Router from '../../router';

export function BaseHeader(props: any) {
  const { Header } = Layout;
  const logoStyle = {
    width: '120px',
    height: '31px',
    background: 'rgba(0, 0, 0, 0.2)',
    margin: '16px 24px 16px 0',
    float: 'left',
  };

  var path = window.location.pathname.split('/');
  var menuKey = ['menu_home'];
  if (path.length >= 2 && path[1].length > 0) {
    menuKey[0] = 'menu_' + path[1];
  }

  return (
    <div>
      <Header>
        <div style={logoStyle} />
        <Menu
          mode="horizontal"
          style={{ lineHeight: '62px', float: 'right' }}
          selectedKeys={menuKey}
        >
          <Menu.Item key="menu_home">
            <Link to={Router.Home.Index}>Home</Link>
          </Menu.Item>
          <Menu.Item key="menu_register">
            <Link to={Router.Home.Register}>注册</Link>
          </Menu.Item>
          <Menu.Item key="menu_login">
            <Link to={Router.Home.Login}>登录</Link>
          </Menu.Item>
        </Menu>
      </Header>
    </div>
  );
}
