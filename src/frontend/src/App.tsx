import "./pages/styles/globals.css";
import "@elastic/eui/dist/eui_theme_light.css";

import React from "react";
import { EuiProvider } from "@elastic/eui";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { SWRConfig } from "swr";

import Notice from "./core/components/notice/Notice";
import { Router } from "./config/router";

import Notfound404 from "./core/components/ui/common/404";

import Home from "./Home";

import User from "./pages/user/User";
import Login from "./pages/auth/Login";
import Register from "./pages/auth/Register";

import ActivateUser from "./pages/auth/ActivateUser";
import RepositoryIndex from "./pages/user/repos/Index";
import RepositoryShow from "./pages/user/repos/Show";
import RepositoryNew from "./pages/user/repos/New";
import Branches from "./pages/user/repos/detail/Branches";
import Detail from "./pages/user/repos/detail/Detail";
import RepositoryLayout from "./pages/user/repos/Layout";

export default function App() {
  // setup();

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
      errorElement: <Notfound404 />,
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
              element: (
                <RepositoryShow defaultChild={<Detail type={"tree"} />} />
              ),
              children: [
                {
                  path: Router.User.Repository.Reference.string(),
                  element: <Detail />,
                },
                {
                  path: Router.User.Repository.Branches.string(),
                  element: <Branches />,
                },
              ],
            },
          ],
        },
      ],
    },
  ]);

  return (
    <>
      <SWRConfig
        value={{ revalidateOnFocus: false, shouldRetryOnError: false }}
      >
        <Notice />
        <EuiProvider colorMode="light">
          <RouterProvider router={router} />
        </EuiProvider>
      </SWRConfig>
    </>
  );
}
