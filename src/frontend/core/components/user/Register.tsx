import React, { useEffect } from 'react';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { Input, Tooltip, Button, Card, Form } from 'antd';
import { withTranslation, WithTranslation } from 'react-i18next';
import { ApolloError, gql } from 'apollo-boost';
import { useMutation } from '@apollo/react-hooks';
import router from 'umi/router';

import { Message } from '../../api/common/notice';
import { UserRules } from '../../api/rule';
import Router from '../../router';

const GQL_REGISTER = gql`
  mutation registerUser($input: NewUserPayload!) {
    registerUser(input: $input) {
      OK
    }
  }
`;

interface NewUserPayload {
  email: string;
  username: string;
  password: string;
}

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 8 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 10 },
  },
};

const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0,
    },
    sm: {
      span: 3,
      offset: 8,
    },
  },
};

function RegisterForm(props: WithTranslation) {
  const { t } = props;

  const [registerUser, { loading: mutationLoading, error: mutationError }] = useMutation<{
    input: NewUserPayload;
  }>(GQL_REGISTER, {
    onCompleted: (data: any) => {
      Message.Success(t('user.tooltip.register_success'));
      router.push(Router.Home.Login);
    },
  });

  const onFinish = function(values: {}) {
    console.info('on finish', values);
    registerUser({
      variables: {
        input: values,
      },
    });
  };

  return (
    <div>
      <Card
        title={t('user.register')}
        extra={<span>{t('user.tooltip.register_notice')}</span>}
        style={{ width: 'auto' }}
      >
        <Form {...formItemLayout} onFinish={onFinish}>
          <Form.Item
            label={t('user.username')}
            name="username"
            rules={[
              {
                required: true,
                message: t('notice.required'),
              },
              {
                min: UserRules.usernameMinLength,
                message: t('user.tooltip.username'),
              },
              {
                max: UserRules.usernameMaxLength,
                message: t('user.tooltip.username'),
              },
              {
                pattern: UserRules.usernameRegex,
                message: t('user.tooltip.username'),
              },
            ]}
          >
            <Input placeholder={t('user.tooltip.username')} />
          </Form.Item>
          <Form.Item
            name="email"
            label={
              <span>
                {t('user.email')}{' '}
                <Tooltip title={t('user.tooltip.email_tip')}>
                  <QuestionCircleOutlined />
                </Tooltip>
              </span>
            }
            rules={[
              {
                required: true,
                message: t('notice.required'),
              },
              {
                type: 'email',
                message: t('user.tooltip.email'),
              },
            ]}
          >
            <Input placeholder={t('user.tooltip.email')} />
          </Form.Item>
          <Form.Item
            name="password"
            label={t('user.password')}
            rules={[
              {
                required: true,
                message: t('notice.required'),
              },
              {
                min: UserRules.pwdMinLength,
                message: t('user.tooltip.password'),
              },
              {
                max: UserRules.pwdMaxLength,
                message: t('user.tooltip.password'),
              },
              {
                pattern: UserRules.passwordRegex,
                message: t('user.tooltip.password'),
              },
            ]}
          >
            <Input.Password placeholder={t('user.tooltip.password')} />
          </Form.Item>
          <Form.Item {...tailFormItemLayout}>
            <Button type="primary" htmlType="submit">
              {t('user.register')}
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
}

export default withTranslation()(RegisterForm);
