import React from "react";
import { EuiTitle } from "@elastic/eui";

interface Props {
  title: string;
}

export default function Title(props: Props) {
  return (
    <EuiTitle size={"s"} className={"mb-5"}>
      <div>{props.title}</div>
    </EuiTitle>
  );
}
