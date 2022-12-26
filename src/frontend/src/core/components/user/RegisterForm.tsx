import React, { useState } from "react";
import validator from "validator";
import { useNavigate } from "react-router-dom";

import {
  EuiButton,
  EuiFieldPassword,
  EuiFieldText,
  EuiForm,
  EuiFormRow,
} from "@elastic/eui";

import { useAuth } from "../../api/auth";
import { Router } from "../../../config/router";
import { useGlobal } from "../../global/global";
import { userRules } from "../../api/rule";
import i18n from "../../i18n/i18n";

export default function RegisterForm(props: any) {
  const auth = useAuth();
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [usernameValidateMsg, setUsernameValidateMsg] = useState(null);
  const [emailValidateMsg, setEmailValidateMsg] = useState(null);
  const [pwdValidateMsg, setPwdValidateMsg] = useState(null);

  const { notice } = useGlobal();
  const navigate = useNavigate();

  const onSubmit = (e: React.MouseEvent) => {
    auth
      .registerUser({
        username: username,
        email: email,
        password: password,
      })
      .then((res) => {
        notice?.success(i18n.t("user.tooltip.register_success"));
        navigate(Router.Home.Login);
        return;
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
      if (
        validator.isEmpty(val) ||
        validator.isLength(val, {
          min: userRules.pwdMinLength,
          max: userRules.pwdMaxLength,
        })
      ) {
        setPwdValidateMsg(i18n.t("user.login_tooltip.password_invalid"));
      } else {
        setPwdValidateMsg(null);
      }
    },
    username: (obj: HTMLInputElement) => {
      const val = obj.value;

      if (
        validator.isEmpty(val) ||
        validator.isLength(val, {
          min: userRules.usernameMinLength,
          max: userRules.usernameMaxLength,
        })
      ) {
        setUsernameValidateMsg(i18n.t("user.login_tooltip.username_invalid"));
      } else {
        setUsernameValidateMsg(null);
      }
    },
  };

  const onBlur = (event: React.FocusEvent<HTMLInputElement>) => {
    const obj = event.target;
    const name = obj.name;
    switch (name) {
      case "username":
        validate.username(obj);
        break;
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
    <div>
      <div className="min-h-full flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <div>
            <h1 className="mx-auto h-12 w-auto text-center text-3xl">
              {i18n.t("user.register") as string}
            </h1>
            <span>{i18n.t<string>("user.tooltip.register_notice")}</span>
          </div>
          <div>
            <div className="-space-y-px shadow-2xl p-8 rounded-xl">
              <EuiForm component="form">
                <EuiFormRow
                  label={i18n.t("user.username") as string}
                  isInvalid={usernameValidateMsg != null}
                  error={usernameValidateMsg}
                >
                  <EuiFieldText
                    name={"username"}
                    required={true}
                    autoFocus={true}
                    onChange={(event) => {
                      setUsername(event.target.value);
                    }}
                    isInvalid={false}
                    onBlur={onBlur}
                  />
                </EuiFormRow>
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
                    onChange={(event): void => {
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
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
