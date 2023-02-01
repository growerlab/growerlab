import React from "react";
import { RepositoryPathGroup } from "../../../common/types";

interface Props extends RepositoryPathGroup {
  Path: string;
}

export default function Blob(props: Props) {
  return <>hello</>;
}
