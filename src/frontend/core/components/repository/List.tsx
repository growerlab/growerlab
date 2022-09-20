import React from "react";
import { useState } from "react";
import Link from "next/link";

import { Router } from "../../../config/router";
import { ListItem } from "./ListItem";
import {
  TypeRepositoriesArgs,
  RepositoryEntity,
} from "../../services/repository/types";
import { Repository } from "../../services/repository/repository";
import { EuiButton, EuiIcon } from "@elastic/eui";

export function RepositoryList(props: TypeRepositoriesArgs) {
  const { ownerPath } = props;
  const [initLoading, setInitLoading] = useState(false);

  const repo = new Repository({ ownerPath: ownerPath });
  const repoData = repo.list();

  if (repoData === null) {
    return (
      <div>
        <div className="text-center text-1xl">
          <EuiIcon type={"folderOpen"} className="text-sky-400 inline" />
          <div className="mt-2">暂无仓库</div>
          <div className="mt-2">
            <Link href={Router.User.Repository.New}>
              <EuiButton type="button" color={"primary"}>
                <EuiIcon type={"plus"}></EuiIcon>
                创建仓库
              </EuiButton>
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const repositories = repoData.repositories;
  // const loadMoreBtn = !initLoading ? (
  //   <div
  //     style={{
  //       textAlign: "center",
  //       marginTop: 12,
  //       height: 32,
  //       lineHeight: "32px",
  //     }}
  //   >
  //     {/* <Button onClick={onLoadMore}>更多</Button> */}
  //   </div>
  // ) : null;

  return (
    <div>
      {repositories.map((repo: RepositoryEntity) => (
        <ListItem repo={repo} />
      ))}
    </div>
  );
}
