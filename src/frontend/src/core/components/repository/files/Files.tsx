import React, { Suspense, useState } from "react";
import { useTitle } from "react-use";
import { useLocation } from "react-router-dom";

import {
  DetailType,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import { getTitle } from "../../../common/document";
import Loading from "../../ui/common/Loading";
import { Tree } from "./Tree";
import Blob from "./Blob";

interface Props extends RepositoryPathGroup {
  type: DetailType;
  blobPath?: string; //
  reference: string;
  initialFolder: string;
  repository?: RepositoryEntity;
}

// for tree and blob
export function Files(props: Props) {
  const { namespace, repo, reference, initialFolder, repository } = props;
  const { blobPath } = props;
  useTitle(getTitle(repo));

  const detail =
    props.type === "blob" ? (
      <Blob {...props} Path={blobPath!} />
    ) : (
      <Tree {...props} />
    );

  return <Suspense fallback={<Loading lines={5} />}>{detail}</Suspense>;
}
