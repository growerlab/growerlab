import React, { useState, Suspense } from "react";
import { useTitle } from "react-use";
import { useNavigate } from "react-router-dom";
import { EuiPanel } from "@elastic/eui";
import { EuiBreadcrumbProps } from "@elastic/eui/src/components/breadcrumbs/breadcrumb";

import {
  DetailType,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import { getTitle } from "../../../common/document";
import { Router } from "../../../../config/router";
import Loading from "../../ui/common/Loading";
import Breadcrumbs from "../../ui/Breadcrumbs";
import { Path } from "../../../common/path";
import Tree from "./Tree";
import Blob from "./Blob";

interface Props extends RepositoryPathGroup {
  type: DetailType;
  blobPath?: string; //
  reference: string;
  initialFolder: string;
  repository?: RepositoryEntity;
}

export function Files(props: Props) {
  const { namespace, repo, reference, initialFolder, repository } = props;
  const { blobPath } = props;
  useTitle(getTitle(repo));

  const navigate = useNavigate();
  const [detailType, setDetailType] = useState(props.type);
  const [isEmptyRepository, setTreeEmpty] = useState<boolean>(false);
  const [currentRepoFolder] = useState(new Path(initialFolder)); // 正在访问的repo路径
  if (!isEmptyRepository && repository?.last_push_at == 0) {
    setTreeEmpty(true);
  }

  const buildTreeURL = (tree: string) => {
    return Router.User.Repository.Reference.render({
      refType: "tree",
      "*": tree,
      ref: reference,
      repo: repo,
    });
  };

  // folder paths
  const folderBreadcrumbs: EuiBreadcrumbProps[] = [
    {
      text: "/",
      onClick: () => {
        const rootURL = buildTreeURL("");
        navigate(rootURL, { state: rootURL });
        currentRepoFolder.reset("");
      },
    },
  ];
  currentRepoFolder.forEach((value, index, array) => {
    if (value === "") {
      return;
    }
    const link = buildTreeURL(value);
    const onClick = () => {
      navigate(link, { state: link });
      setDetailType("tree");
      currentRepoFolder.reset(array.slice(0, index + 1));
    };
    folderBreadcrumbs.push({
      text: value,
      onClick: index == array.length - 1 ? undefined : onClick,
    });
  });

  const detail =
    detailType === "blob" ? (
      <Blob {...props} Path={blobPath!} />
    ) : (
      <Tree
        {...props}
        initialFolder={currentRepoFolder.toString()}
        onChange={(name, type) => {
          // TODO ...
          return;
        }}
      />
    );

  return (
    <Suspense fallback={<Loading lines={5} />}>
      <EuiPanel hasBorder={true} paddingSize={"s"}>
        <Breadcrumbs
          truncate={true}
          max={4}
          breadcrumbs={folderBreadcrumbs}
          icon={"submodule"}
        />
        {detail}
      </EuiPanel>
    </Suspense>
  );
}
