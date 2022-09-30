import React from "react";
import { Link } from "react-router-dom";

import { Router } from "../../../config/router";
import { Item } from "./Item";
import {
  RepositoriesArgs,
  TypeRepository,
} from "../../services/repository/types";
import { Repository } from "../../services/repository/repository";
import { EuiButton, EuiIcon, EuiEmptyPrompt } from "@elastic/eui";
import { useGlobal } from "../../global/init";

export function RepositoryList(props: RepositoriesArgs) {
  const { ownerPath } = props;
  // const [initLoading, setInitLoading] = useState(false);
  const global = useGlobal();

  const repo = new Repository(ownerPath);
  const repoData = repo.list();

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
    <div>
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
