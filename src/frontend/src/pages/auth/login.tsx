import React from "react";
import { getTitle } from "../../core/common/document";
import LoginForm from "../../core/components/user/Login";
import i18n from "../../core/i18n/i18n";

export default function Login(props: any) {
  return (
    <div>
      <div>{getTitle(i18n.t("website.login"))}</div>
      <LoginForm />
    </div>
  );
}
