import React, { ChangeEvent, useState } from "react";
import validator from "validator";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiFieldPassword,
  EuiButton,
} from "@elastic/eui";
import { redirect } from "react-router-dom";

import { Router } from "../../../config/router";
import { useGlobal } from "../../global/init";
import { Session } from "../../services/auth/session";
import { Auth } from "../../services/auth/auth";

export default function LoginForm(props: any) {
  const global = useGlobal();
  const notice = global.notice!;

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailValidateMsg, setEmailValidateMsg] = useState(null);
  const [pwdValidateMsg, setPwdValidateMsg] = useState(null);

  const { t } = props;

  const onSubmit = (e: React.MouseEvent) => {
    const service = new Auth();
    service.login(email, password).then((res) => {
      if (res === undefined) {
        notice.error(t("user.tooltip.login_fail"));
        return;
      }

      Session.storeLogin(res.data);
      notice.success(t("user.tooltip.login_success"));
      redirect(Router.User.Index);
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
