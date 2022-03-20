import React from 'react';
import {withTranslation} from 'react-i18next';
import {getTitle} from '../common/document';
import Head from 'next/head'
import Link from 'next/link'
import {Button} from "evergreen-ui";
import {Router} from "../config/router";


const index = function (props: any) {
  return (
    <div>
      <Head>
        <title>{getTitle()}</title>
      </Head>

      <h2 className="text-6xl font-bold text-center mt-7">Rethinking Git</h2>
      <div className="text-center mt-7">
        <Link href={Router.Home.Login}>
          <Button marginRight={16} appearance="primary" size={"large"}>
            Login
          </Button>
        </Link>
      </div>
    </div>
  );
};

export default withTranslation()(index);
