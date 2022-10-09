import React from "react";
import { useTitle } from "react-use";
import { EuiButton } from "@elastic/eui";

import { Router } from "./config/router";
import HomeLayout from "./pages/layouts/HomeLayout";
import { getTitle } from "./core/common/document";

export default function Home() {
  useTitle(getTitle());

  return (
    <div>
      <HomeLayout>
        <h2 className="text-6xl font-bold text-center mt-7">Rethinking Git</h2>
        <div className="text-center mt-7">
          <EuiButton href={Router.Home.Login}>Login</EuiButton>
        </div>
      </HomeLayout>
    </div>
  );
}
