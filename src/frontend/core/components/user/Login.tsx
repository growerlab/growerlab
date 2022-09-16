import React, { ChangeEvent, useState } from "react";
import { withTranslation, WithTranslation } from "react-i18next";
import { useRouter } from "next/router";
import validator from "validator";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiFieldPassword,
  EuiButton,
} from "@elastic/eui";

import { Router } from "../../../config/router";
import { LoginService } from "../../services/auth/login";
import { useGlobal } from "../../global/init";

function LoginForm(props: WithTranslation) {
  const global = useGlobal();
  const notice = global.notice!;
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailValidateMsg, setEmailValidateMsg] = useState(null);
  const [pwdValidateMsg, setPwdValidateMsg] = useState(null);

  const { t } = props;

  const onSubmit = (e: React.MouseEvent) => {
    const service = new LoginService();
    service.login(email, password).then((res) => {
      if (res === undefined) {
        notice.error(t("user.tooltip.login_fail"));
        return;
      }
      notice.success(t("user.tooltip.login_success"));
      router.push(Router.User.Index);
    });
  };

  const validate = {
    email: (obj: HTMLInputElement) => {
      const val = obj.value;
      if (!validator.isEmail(val)) {
        setEmailValidateMsg(t("user.login_tooltip.email_invalid"));
      } else {
        setEmailValidateMsg(null);
      }
    },
    password: (obj: HTMLInputElement) => {
      const val = obj.value;
      if (validator.isEmpty(val)) {
        setPwdValidateMsg(t("user.login_tooltip.password_invalid"));
      } else {
        setPwdValidateMsg(null);
      }
    },
  };

  const onBlur = (event: React.FocusEvent<HTMLInputElement>) => {
    const obj = event.target;
    const type = obj.type;
    switch (type) {
      case "email":
        validate.email(obj);
        break;
      case "password":
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
            <EuiForm component="form">
              <EuiFormRow
                label="Email"
                isInvalid={emailValidateMsg != null}
                error={emailValidateMsg}
              >
                <EuiFieldText
                  name={"email"}
                  type={"email"}
                  required={true}
                  onChange={(event) => {
                    setEmail(event.target.value);
                  }}
                  isInvalid={false}
                  onBlur={onBlur}
                />
              </EuiFormRow>
              <EuiFormRow
                label="Password"
                isInvalid={pwdValidateMsg != null}
                error={pwdValidateMsg}
              >
                <EuiFieldPassword
                  type={"password"}
                  name={"password"}
                  isInvalid={false}
                  required={true}
                  onBlur={onBlur}
                  onChange={(event: ChangeEvent<HTMLInputElement>): void => {
                    setPassword(event.target.value);
                  }}
                />
              </EuiFormRow>
              <EuiFormRow>
                <EuiButton
                  fill
                  color="primary"
                  onClick={onSubmit}
                  className={"w-full"}
                >
                  {t("user.login")}
                </EuiButton>
              </EuiFormRow>
            </EuiForm>
          </div>
        </div>
      </div>
    </div>
  );
}

export default withTranslation()(LoginForm);
