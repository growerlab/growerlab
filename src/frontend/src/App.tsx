import "./pages/styles/globals.css";

import "@elastic/eui/dist/eui_theme_light.css";

import React from "react";
import { RecoilRoot } from "recoil";
import { EuiProvider } from "@elastic/eui";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import Notice from "./core/components/notice/Notice";
import { setup } from "./core/global/init";
import Home from "./Home";
import User from "./pages/user/User";
import Login from "./pages/auth/Login";

export default function App() {
  setup();

  const router = createBrowserRouter([
    {
      path: "/",
      element: <Home />,
    },
    {
      path: "/login",
      element: <Login />,
    },
    {
      path: "/user",
      element: <User />,
    },
  ]);

  return (
    <div>
      <RecoilRoot>
        <Notice />
        <EuiProvider colorMode="light">
          <RouterProvider router={router} />
        </EuiProvider>
      </RecoilRoot>
    </div>
  );
}
