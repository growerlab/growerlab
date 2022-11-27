import "./pages/styles/globals.css";
import "@elastic/eui/dist/eui_theme_light.css";

import React from "react";
import { RecoilRoot } from "recoil";
import { EuiProvider } from "@elastic/eui";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import Notice from "./core/components/notice/Notice";
import { setup } from "./core/global/init";
import { Router } from "./config/router";

import Home from "./Home";
import User from "./pages/user/User";
import Login from "./pages/auth/Login";
import Register from "./pages/auth/Register";
import ActivateUser from "./pages/auth/ActivateUser";
import RepositoryIndex from "./pages/user/repos/Index";
import RepositoryShow from "./pages/user/repos/Show";
import RepositoryNew from "./pages/user/repos/New";
import Error404 from "./pages/common/404";
import Branchs from "./pages/user/repos/detail/Branchs";
import Detail from "./pages/user/repos/detail/Detail";
import RepositoryLayout from "./pages/user/repos/Layout";

export default function App() {
  setup();

  const router = createBrowserRouter([
    {
      path: Router.Home.Index,
      element: <Home />,
    },
    {
      path: Router.Home.Login,
      element: <Login />,
    },
    {
      path: Router.Home.Register,
      element: <Register />,
    },
    {
      path: Router.Home.ActivateUser,
      element: <ActivateUser />,
    },
    {
      path: Router.User.Index,
      element: <User />,
      errorElement: <Error404 />,
      children: [
        {
          path: Router.User.Repository.Index,
          element: <RepositoryLayout defaultChild={<RepositoryIndex />} />,
          children: [
            {
              path: Router.User.Repository.New,
              element: <RepositoryNew />,
            },
            {
              path: Router.User.Repository.Show.string(),
              element: <RepositoryShow defaultChild={<Detail />} />,
              children: [
                {
                  path: Router.User.Repository.Branchs.string(),
                  element: <Branchs />,
                },
              ],
            },
          ],
        },
      ],
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
