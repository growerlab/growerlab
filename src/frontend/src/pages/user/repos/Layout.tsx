import React from "react";
import { useOutlet } from "react-router-dom";

interface Props extends React.PropsWithChildren {
  defaultChild?: React.ReactElement;
}

export default function RepositoryLayout(props: Props) {
  const outlet = useOutlet();
  const defaultOutlet = outlet === null ? props.defaultChild : outlet;

  return <div>{defaultOutlet}</div>;
}
