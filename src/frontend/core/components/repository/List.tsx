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
import {
  EuiButton,
  EuiIcon,
  EuiButtonEmpty,
  EuiEmptyPrompt,
} from "@elastic/eui";

export function RepositoryList(props: TypeRepositoriesArgs) {
  const { ownerPath } = props;
  // const [initLoading, setInitLoading] = useState(false);

  const repo = new Repository({ ownerPath: ownerPath });
  const repoData = repo.list();

  if (repoData === null) {
    return (
      <div>
        <EuiEmptyPrompt
          title={<h2>无任何仓库，立即创建！</h2>}
          actions={[
            <Link href={Router.User.Repository.New}>
              <EuiButton type="button" color={"primary"}>
                <EuiIcon type={"plus"}></EuiIcon>
                创建仓库
              </EuiButton>
            </Link>,
            // <EuiButtonEmpty color="primary">Start a trial</EuiButtonEmpty>,
          ]}
        />
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
        <ListItem repo={repo} key={repo.uuid} />
      ))}
    </div>
  );
}
