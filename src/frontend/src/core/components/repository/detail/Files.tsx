import React from "react";
import { RepositoryPathGroup } from "../../../common/types";
import Loading from "../../common/Loading";

interface Props extends RepositoryPathGroup {
  branch?: "master" | string;
}

export function Files(props: Props) {
  return (
    <>
      <Loading lines={5}></Loading>
    </>
  );
}
