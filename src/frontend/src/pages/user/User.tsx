import React from "react";
import { Outlet } from "react-router-dom";

import UserLayout from "../layouts/UserLayout";

export default function User() {
  return (
    <div>
      <UserLayout title="Dashboard">
        <Outlet />
      </UserLayout>
    </div>
  );
}
