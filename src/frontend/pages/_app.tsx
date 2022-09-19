import "../styles/globals.css";

import React from "react";
import type { AppProps } from "next/app";
import { RecoilRoot } from "recoil";
import { EuiProvider } from "@elastic/eui";

import "@elastic/eui/dist/eui_theme_light.css";

import Notice from "../core/components/notice/Notice";
import NoSSR from "../core/components/global/NoSSR";
import { setup } from "../core/global/init";

function MyApp({ Component, pageProps }: AppProps) {
  setup();

  return (
    <NoSSR>
      <RecoilRoot>
        <Notice />
        <EuiProvider colorMode="light">
          <Component {...pageProps} />
        </EuiProvider>
      </RecoilRoot>
    </NoSSR>
  );
}

export default MyApp;
