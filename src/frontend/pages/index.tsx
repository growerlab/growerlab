import React, { useContext, useEffect } from "react";
import { withTranslation } from "react-i18next";
import Head from "next/head";
import { getTitle } from "../core/common/document";

import { EuiButton } from "@elastic/eui";
import { Router } from "../config/router";

const index = function (props: any) {
  return (
    <div>
      <Head>
        <title>{getTitle()}</title>
      </Head>

      <h2 className="text-6xl font-bold text-center mt-7">Rethinking Git</h2>
      <div className="text-center mt-7">
        <EuiButton href={Router.Home.Login}>Login</EuiButton>
      </div>
    </div>
  );
};

export default withTranslation()(index);
