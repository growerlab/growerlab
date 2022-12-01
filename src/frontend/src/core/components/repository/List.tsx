import React from "react";
import { Link } from "react-router-dom";

import { Router } from "../../../config/router";
import { Item } from "./Item";
import { RepositoriesNamespace, TypeRepository } from "../../common/types";
import { useRepositoryAPI } from "../../api/repository/repository";
import { EuiButton, EuiIcon, EuiEmptyPrompt } from "@elastic/eui";
import { useGlobal } from "../../global/init";

export function RepositoryList(props: RepositoriesNamespace) {
  const { namespace } = props;
  const repositoryAPI = useRepositoryAPI(namespace);
  const global = useGlobal();

  const repoData = repositoryAPI.list();

  if (repoData === null) {
    return (
      <div>
        <EuiEmptyPrompt
          title={<h2>无任何仓库，立即创建！</h2>}
          actions={[
            <Link to={Router.User.Repository.New} key={""}>
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

  const repositories = repoData;
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
    <div className="p-5">
      {repositories?.map((repo: TypeRepository) => (
        <Item
          global={global}
          repo={repo.repository}
          key={repo.repository.uuid}
        />
      ))}
    </div>
  );
}
