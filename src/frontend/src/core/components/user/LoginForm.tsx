import React, { ChangeEvent, useState } from "react";
import validator from "validator";
import {
  EuiFieldText,
  EuiForm,
  EuiFormRow,
  EuiFieldPassword,
  EuiButton,
} from "@elastic/eui";

import i18n from "../../i18n/i18n";
import { useSession, useNotice, useAuth } from "../../global/state";

interface Options {
  onSuccess: () => any;
}

export default function LoginForm(props: Options) {
  const notice = useNotice();
  const { onSuccess } = props;
  const auth = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [emailValidateMsg, setEmailValidateMsg] = useState(null);
  const [pwdValidateMsg, setPwdValidateMsg] = useState(null);
  const session = useSession();

  const onSubmit = (e: React.MouseEvent) => {
    auth.login(email, password).then((res) => {
      if (res === undefined) {
        notice.error(i18n.t("user.tooltip.login_fail"));
        return;
      }

      session.storeLogin(res.data);
      notice.success(i18n.t("user.tooltip.login_success"));
      onSuccess();
    });
  };

  const validate = {
    email: (obj: HTMLInputElement) => {
      const val = obj.value;
      if (!validator.isEmail(val)) {
        setEmailValidateMsg(i18n.t("user.login_tooltip.email_invalid"));
      } else {
        setEmailValidateMsg(null);
      }
    },
    password: (obj: HTMLInputElement) => {
      const val = obj.value;
      if (validator.isEmpty(val)) {
        setPwdValidateMsg(i18n.t("user.login_tooltip.password_invalid"));
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
    <>
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
            icon={"user"}
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
            {i18n.t<string>("user.login")}
          </EuiButton>
        </EuiFormRow>
      </EuiForm>
    </>
  );
}
