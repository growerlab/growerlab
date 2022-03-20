import React, { useState, useEffect } from 'react';
import { WithTranslation, withTranslation } from 'react-i18next';
import {
  FrownTwoTone,
  QuestionCircleTwoTone,
  SmileTwoTone,
  LoadingOutlined,
} from '@ant-design/icons';
import { Result, Button } from 'antd';
import { gql } from 'apollo-boost';
import { useMutation } from '@apollo/react-hooks';
import router from 'umi/router';
import { useParams } from 'react-router';
import { Router } from '../../router';

const GQL_REGISTER = gql`
  mutation activateUser($input: ActivationCodePayload!) {
    activateUser(input: $input) {
      OK
    }
  }
`;

interface ActivationCodePayload {
  code: string;
}

interface Status {
  Title: string;
  Status?: string;
  SubTitle?: string;
  Icon?: React.ReactNode;
  Extra?: React.ReactNode;
}

// 状态
//  1. 请求参数中未包含code 2. 请求接口中 3|4. 接口返回正常|错误  5. 激活码已被使用过(服务器端返回)
//
function Activate(props: WithTranslation) {
  const { t } = props;

  let { code } = useParams();

  let loginBtn = (
    <Button
      type="primary"
      onClick={() => {
        router.push(Router.Home.Login);
      }}
    >
      {t('user.login')}
    </Button>
  );

  const status: { [key: string]: Status } = {
    NotFound: {
      Title: t('user.activate.not_found.code'),
      SubTitle: t('user.activate.invalid'),
      Icon: <QuestionCircleTwoTone />,
    },
    Pending: {
      Title: t('user.activate.pending'),
      SubTitle: t('user.activate.pending_sub'),
      Icon: <LoadingOutlined />,
    },
    Failed: {
      Title: t('user.activate.invalid'),
      Icon: <FrownTwoTone />,
      Extra: loginBtn,
    },
    Success: {
      Title: t('user.activate.success'),
      SubTitle: t('user.activate.success_sub'),
      Icon: <SmileTwoTone />,
      Extra: loginBtn,
    },
  };

  const [curStatus, setStatus] = useState(status['Pending']);

  const [activateUser] = useMutation<{
    input: ActivationCodePayload;
  }>(GQL_REGISTER);

  useEffect(() => {
    activateUser({
      variables: {
        input: {
          code: code,
        },
      },
    })
      .then((data: any) => {
        if (data.data.activateUser.OK !== null) {
          setStatus(status['Success']);
        }
      })
      .catch((reason: any) => {
        setStatus(status['Failed']);
      });
  }, []);

  return (
    <Result
      // status={curStatus.Status !== null ? 'warning' : curStatus.Status}
      icon={curStatus.Icon === null ? '' : curStatus.Icon}
      title={curStatus.Title}
      subTitle={curStatus.SubTitle}
      extra={curStatus.Extra}
    />
  );
}

export default withTranslation()(Activate);
