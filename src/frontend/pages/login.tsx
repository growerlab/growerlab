import React from 'react';
import LoginForm from '../components/user/Login';
import {withTranslation} from 'react-i18next';
import {getTitle} from '../common/document';
import Head from 'next/head'

const login = function (props: any) {
  const {t} = props;

  return (
    <div>
      <Head>
        <title>{getTitle(t('website.login'))}</title>
      </Head>
      <LoginForm/>
    </div>
  );
};

export default withTranslation()(login);
