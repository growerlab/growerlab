import React from "react";
import { useOutlet } from "react-router-dom";

import UserLayout from "../layouts/UserLayout";

interface Props extends React.PropsWithChildren {
  defaultChild?: React.ReactElement;
}

export default function User(props: Props) {
  const outlet = useOutlet();
  const defaultOutlet = outlet === null ? props.defaultChild : outlet;

  return (
    <div>
      <UserLayout title="Dashboard">{defaultOutlet}</UserLayout>
    </div>
  );
}
