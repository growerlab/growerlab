import "./pages/styles/globals.css";

import { withTranslation } from "react-i18next";
import React from "react";
import { RecoilRoot } from "recoil";
import { EuiProvider } from "@elastic/eui";
import { Route, Routes } from "react-router-dom";

import Notice from "./core/components/notice/Notice";
import { setup } from "./core/global/init";

import "@elastic/eui/dist/eui_theme_light.css";
import Home from "./Home";
import Error404 from "./pages/common/404";

function App() {
  setup();

  return (
    <div>
      <RecoilRoot>
        <Notice />
        <EuiProvider colorMode="light">
          <Routes>
            <Route path="/user">
              <App />
            </Route>
            <Route path="/">
              <Home />
            </Route>
            <Route>
              <Error404 />
            </Route>
          </Routes>
        </EuiProvider>
      </RecoilRoot>
    </div>
  );
}

export default withTranslation()(App);
