import React from "react";
import { useTitle } from "react-use";

import { getTitle } from "../../core/common/document";
import LoginForm from "../../core/components/user/LoginForm";
import i18n from "../../core/i18n/i18n";

export default function Login(props: any) {
  useTitle(getTitle(i18n.t("website.login")));

  return (
    <>
      <LoginForm />
    </>
  );
}
