import React from "react";
import { EuiButton } from "@elastic/eui";

import { Router } from "./config/router";

const Home = () => {
  return (
    <div>
      {/* <Head>
        <title>{getTitle()}</title>
      </Head> */}

      <h2 className="text-6xl font-bold text-center mt-7">Rethinking Git</h2>
      <div className="text-center mt-7">
        <EuiButton href={Router.Home.Login}>Login</EuiButton>
      </div>
    </div>
  );
};

export default Home;
