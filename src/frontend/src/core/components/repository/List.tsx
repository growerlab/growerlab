import React from "react";
import { Link } from "react-router-dom";
import useSWRImmutable, { Fetcher } from "swr";
import {
  EuiButton,
  EuiIcon,
  EuiEmptyPrompt,
  EuiHorizontalRule,
} from "@elastic/eui";

import { Router } from "../../../config/router";
import { Header } from "./Header";
import {
  RepositoriesNamespace,
  RepositoryEntity,
  TypeRepositories,
} from "../../common/types";
import { useRepositoryAPI } from "../../api/repository";
import { useGlobal } from "../../global/global";
import i18n from "../../i18n/i18n";

export function RepositoryList(props: RepositoriesNamespace) {
  const { namespace } = props;
  const repositoryAPI = useRepositoryAPI(namespace);
  const global = useGlobal();

  const fetcher: Fetcher<TypeRepositories> = () => {
    return repositoryAPI.list().then((res) => {
      return res.data;
    });
  };
  const { data } = useSWRImmutable<TypeRepositories>(
    `/swr/key/repos/${namespace}`,
    fetcher
  );

  if (data?.repositories.length == 0) {
    return (
      <div>
        <EuiEmptyPrompt
          title={
            <h2>{i18n.t<string>("repository.tooltip.no_repositories")}</h2>
          }
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
        <>
          <Header global={global} repo={repo} />
          <EuiHorizontalRule />
        </>
      ))}
    </div>
  );
}
