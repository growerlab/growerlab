import React from "react";
import { withTranslation } from "react-i18next";
import { getTitle } from "../../core/common/document";
import LoginForm from "../../core/components/user/Login";

const login = function (props: any) {
  const { t } = props;

  return (
    <div>
      <div>{getTitle(t("website.login"))}</div>
      <LoginForm />
    </div>
  );
};

export default withTranslation()(login);
