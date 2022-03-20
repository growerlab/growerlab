import React, { ChangeEvent, useEffect, useState } from 'react';
import { withTranslation, WithTranslation } from 'react-i18next';
import { useRouter } from 'next/router';
import Link from 'next/link';
import { TextInputField, Button, SearchIcon } from 'evergreen-ui';
import validator from 'validator';

import { Message } from '../../api/common/notice';
import { UserRules } from '../../api/rule';
import { Router } from '../../config/router';
import { LoginService } from '../../services/auth/login';

interface LoginUserPayload {
  email: string;
  password: string;
}

function LoginForm(props: WithTranslation) {
  const router = useRouter();
  const { t } = props;
  const [emailValidateMsg, setEmailValidateMsg] = useState(null);
  const [pwdValidateMsg, setPwdValidateMsg] = useState(null);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const onSubmit = (e: React.MouseEvent) => {
    LoginService.login(email, password).then((res) => {
      if (res === undefined) {
        Message.Error(t('user.tooltip.login_fail'));
        return;
      }
      Message.Success(t('user.tooltip.login_success'));
      router.push(Router.User.Index);
    });
  };

  const validate = {
    email: (obj: HTMLInputElement) => {
      const val = obj.value;
      if (!validator.isEmail(val)) {
        setEmailValidateMsg(t('user.login_tooltip.email_invalid'));
      } else {
        setEmailValidateMsg(null);
      }
    },
    password: (obj: HTMLInputElement) => {
      const val = obj.value;
      if (validator.isEmpty(val)) {
        setPwdValidateMsg(t('user.login_tooltip.password_invalid'));
      } else {
        setPwdValidateMsg(null);
      }
    },
  };

  const onBlur = (event: React.FocusEvent<HTMLInputElement>) => {
    const obj = event.target;
    const type = obj.type;
    switch (type) {
      case 'email':
        validate.email(obj);
        break;
      case 'password':
        validate.password(obj);
        break;
    }
    return;
  };

  return (
    <div className="min-h-full flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h1 className="mx-auto h-12 w-auto text-center text-3xl">Logo</h1>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to your account
          </h2>
        </div>
        <div>
          <div className="-space-y-px shadow-2xl p-8 rounded-xl">
            <div>
              <TextInputField
                type="email"
                name="email"
                autoComplete="email"
                isInvalid={false}
                required
                label="Email"
                validationMessage={emailValidateMsg}
                onBlur={onBlur}
                onChange={(event: ChangeEvent<HTMLInputElement>): void =>
                  setEmail(event.target.value)
                }
              />
            </div>
            <div>
              <TextInputField
                type="password"
                name="password"
                isInvalid={false}
                required
                label="Password"
                validationMessage={pwdValidateMsg}
                onBlur={onBlur}
                onChange={(event: ChangeEvent<HTMLInputElement>): void =>
                  setPassword(event.target.value)
                }
              />
            </div>
            <div>
              <Button
                appearance="primary"
                marginY={8}
                marginRight={12}
                className="w-full"
                size="medium"
                onClick={onSubmit}
              >
                {t('user.login')}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default withTranslation()(LoginForm);
