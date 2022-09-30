import React from "react";
import Head from "next/head";
import { withTranslation } from "react-i18next";

import { getTitle } from "../core/common/document";
import LoginForm from "../core/components/user/Login";

const login = function (props: any) {
  const { t } = props;

  return (
    <div>
      <Head>
        <title>{getTitle(t("website.login"))}</title>
      </Head>
      <LoginForm />
    </div>
  );
};

export default withTranslation()(login);
