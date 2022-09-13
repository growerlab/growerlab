import "../styles/globals.css";

import React from "react";
import type { AppProps } from "next/app";
import { RecoilRoot } from "recoil";
import { EuiProvider } from "@elastic/eui";

import "@elastic/eui/dist/eui_theme_light.css";

import Notice from "../core/components/notice/Notice";
import { setup } from "../core/global/init";

function MyApp({ Component, pageProps }: AppProps) {
  setup();

  return (
    <RecoilRoot>
      <Notice />
      <EuiProvider>
        <Component {...pageProps} />
      </EuiProvider>
    </RecoilRoot>
  );
}

export default MyApp;
