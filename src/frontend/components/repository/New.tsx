import {Form, Button, Input, Checkbox} from 'antd';
import React from 'react';
import {gql} from 'apollo-boost';
import {useMutation} from '@apollo/react-hooks';
import {RepositoryRules} from '../../api/rule';
import {Message} from '../../api/common/notice';
import {getUserInfo} from '../../api/auth/session';
import {Router} from '../../config/router';
import i18n from '../../i18n/i18n';

const GQL_REGISTER = gql`
  mutation createRepository($input: NewRepositoryPayload!) {
    createRepository(input: $input) {
      OK
    }
  }
`;

interface NewRepositoryPayload {
  namespacePath: string;
  name: string;
  public: boolean;
}

export function NewRepositoryFrom(props: any) {
  const {t} = props;
  const [form] = Form.useForm();

  const [createRepository, {loading: mutationLoading, error: mutationError}] = useMutation<{
    input: NewRepositoryPayload;
  }>(GQL_REGISTER, {
    onCompleted: (data: any) => {
      Message.Success(t('repository.tooltip.success'));
      router.push(Router.User.Repository.List);
    },
  });

  const handleSubmit = (values: Store) => {
    var payload = values as NewRepositoryPayload;

    createRepository({
      variables: {
        input: payload,
      },
    });
  };

  const formItemLayout = {
    labelCol: {span: 4},
    wrapperCol: {span: 14},
  };

  return (
    <Form {...formItemLayout} onFinish={handleSubmit}>
      <Form.Item
        label={i18n.t('repository.owner')}
        name="namespacePath"
        initialValue={getUserInfo()!.namespacePath}
      >
        <Input hidden={true}/>
        <span className="ant-form-text">{getUserInfo()!.name}</span>
      </Form.Item>

      <Form.Item
        label={i18n.t('repository.name')}
        name="name"
        rules={[
          {
            required: true,
            message: i18n.t('notice.required'),
          },
          {
            pattern: RepositoryRules.repositoryNameRegex,
            message: i18n.t('repository.tooltip.name'),
          },
        ]}
      >
        <Input placeholder={i18n.t('repository.name')}/>
      </Form.Item>

      <Form.Item
        label={i18n.t('repository.description')}
        name="description"
        rules={[
          {
            required: false,
          },
        ]}
      >
        <Input/>
      </Form.Item>

      <Form.Item name="public" label={i18n.t('repository.public')} initialValue={true}>
        <Checkbox checked={true}/>
      </Form.Item>

      <Form.Item wrapperCol={{span: 6, offset: 4}}>
        <Button type="primary" htmlType="submit">
          {i18n.t('repository.create_repository')}
        </Button>
      </Form.Item>
    </Form>
  );
}
