import React from "react";
import { Link } from "react-router-dom";
import useSWR, { Fetcher } from "swr";
import { EuiButton, EuiIcon, EuiEmptyPrompt } from "@elastic/eui";

import { Router } from "../../../config/router";
import { Item } from "./Item";
import {
  RepositoriesNamespace,
  RepositoryEntity,
  TypeRepositories,
} from "../../common/types";
import { useRepositoryAPI } from "../../api/repository/repository";
import { useGlobal } from "../../global/global";

export function RepositoryList(props: RepositoriesNamespace) {
  const { namespace } = props;
  const repositoryAPI = useRepositoryAPI(namespace);
  const global = useGlobal();

  const fetcher: Fetcher<TypeRepositories> = () => {
    return repositoryAPI.list().then((res) => {
      return res.data;
    });
  };
  const { data, error } = useSWR<TypeRepositories>(
    `/swr/key/repos/${namespace}`,
    fetcher
  );

  if (
    error === undefined ||
    data === undefined ||
    data?.repositories.length == 0
  ) {
    if (error === null) {
      console.error(error);
    }
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

  return (
    <div className="p-5">
      {data?.repositories.map((repo: RepositoryEntity) => (
        <Item global={global} repo={repo} key={repo.uuid} />
      ))}
    </div>
  );
}
