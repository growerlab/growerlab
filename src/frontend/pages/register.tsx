import React from 'react';
import {Row, Col} from 'antd';
import RegisterForm from '../components/user/Register';
import {withTranslation} from 'react-i18next';
import {getTitle} from '../common/document';

const register = function Register(props: any) {
  const {t} = props;
  getTitle(t('website.register'));

  return (
    <div>
      <Row>
        <Col span={6}/>
        <Col span={12}>
          <RegisterForm/>
        </Col>
        <Col span={6}/>
      </Row>
    </div>
  );
};

export default withTranslation()(register);
