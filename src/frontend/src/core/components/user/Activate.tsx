import React, { useState, useEffect } from "react";
import { WithTranslation, withTranslation } from "react-i18next";
import { EuiButton, EuiIcon, EuiEmptyPrompt } from "@elastic/eui";

import { useParams, redirect } from "react-router-dom";
import { Router } from "../../../config/router";
import { Auth } from "../../services/auth/auth";
import Error404 from "../../../pages/common/404";

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

  const loginBtn = (
    <EuiButton
      color="primary"
      onClick={() => {
        redirect(Router.Home.Login);
      }}
    >
      {t("user.login")}
    </EuiButton>
  );

  const status: { [key: string]: Status } = {
    NotFound: {
      Title: t("user.activate.not_found.code"),
      SubTitle: t("user.activate.invalid"),
      Icon: <EuiIcon type="alert" size={"xl"} color={"warning"} />,
    },
    Pending: {
      Title: t("user.activate.pending"),
      SubTitle: t("user.activate.pending_sub"),
      Icon: <EuiIcon type="alert" size={"xl"} color={"warning"} />,
    },
    Failed: {
      Title: t("user.activate.invalid"),
      Icon: <EuiIcon type="alert" size={"xl"} color={"warning"} />,
      Extra: loginBtn,
    },
    Success: {
      Title: t("user.activate.success"),
      SubTitle: t("user.activate.success_sub"),
      Icon: <EuiIcon type="check" size={"xl"} color={"success"} />,
      Extra: loginBtn,
    },
  };

  const [curStatus, setStatus] = useState(status["Pending"]);
  const { code } = useParams();
  if (code === undefined) {
    return <Error404 />;
  }

  const auth = new Auth();
  auth
    .activate(code)
    .then(() => {
      setStatus(status["Success"]);
    })
    .catch((reason: any) => {
      setStatus(status["Failed"]);
    });

  return (
    <>
      <EuiEmptyPrompt
        iconType="securityAnalyticsApp"
        iconColor="default"
        title={<h2>{curStatus.Title}</h2>}
        titleSize="xs"
        body={<p>{curStatus.SubTitle}</p>}
        actions={curStatus.Extra}
      />
    </>
  );
}

export default withTranslation()(Activate);
