import React, { useState, Suspense } from "react";
import useSWR from "swr";
import { EuiBasicTable, EuiIcon, EuiLink, EuiPanel } from "@elastic/eui";
import { useTitle } from "react-use";
import TimeAgo from "timeago-react";
import { Link } from "react-router-dom";

import {
  FileEntity,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import Loading from "../../common/Loading";
import { useRepositoryAPI } from "../../../api/repository";
import EmptyTree from "./EmptyTree";
import { getTitle } from "../../../common/document";
import i18n from "../../../i18n/i18n";
import { Router } from "../../../../config/router";

interface Props extends RepositoryPathGroup {
  reference: string;
  folder: string;
  repository?: RepositoryEntity;
}

export function Files(props: Props) {
  const { namespace, repo, reference, folder, repository } = props;

  useTitle(getTitle(repo));

  const repositoryAPI = useRepositoryAPI(namespace);
  const [isEmptyTree, setTreeEmpty] = useState<boolean>(false);
  const [currentRepoFolder, setCurrentRepoFolder] = useState(folder); // 正在访问的repo路径

  if (!isEmptyTree && repository?.last_push_at == 0) {
    setTreeEmpty(true);
  }

  // api
  const fetcher = () => {
    const params = {
      namespace,
      repo,
      ref: reference,
      folder: currentRepoFolder,
    };
    return repositoryAPI.treeFiles(params).then((res) => {
      return res.data;
    });
  };
  const { data } = useSWR<FileEntity[]>(
    isEmptyTree
      ? null
      : `/swr/key/repo/${namespace}/${repo}/${window.location.pathname}`,
    fetcher,
    { suspense: true }
  );

  if (isEmptyTree) {
    if (repository === undefined) {
      return <></>;
    }
    return (
      <EmptyTree
        cloneURLSSH={repository.git_ssh_url}
        cloneURLHttp={repository.git_http_url}
        defaultBranch={repository.default_branch}
      />
    );
  }

  const columns: any = [
    {
      field: "name",
      name: i18n.t("repository.commit.file_name"),
      render: (name: string, record: FileEntity) => {
        const icon = record.is_file ? (
          <EuiIcon type={"document"} />
        ) : (
          <EuiIcon type={"folderClosed"} />
        );
        const folderPath = folder != "" ? folder + "/" + name : name;
        const link = record.is_file
          ? Router.User.Repository.Blob.render({
              filepath: name,
              ref: reference,
              repo: repo,
            })
          : Router.User.Repository.Tree.render({
              "*": folderPath,
              ref: reference,
              repo: repo,
            });
        return (
          <>
            {icon}{" "}
            <Link to={link} onClick={() => setCurrentRepoFolder(folderPath)}>
              {name}
            </Link>
          </>
        );
      },
    },
    {
      field: "last_commit_message",
      name: i18n.t("repository.commit.lastCommitMessage"),
      truncateText: true,
      render: (last_commit_message: string) => (
        <EuiLink href="#commit详情">{last_commit_message}</EuiLink>
      ),
    },
    {
      field: "last_commit_date",
      name: i18n.t("repository.commit.lastCommitDate"),
      dataType: "date",
      align: "right",
      render: (last_commit_date: number) => {
        const lastDate = new Date(last_commit_date * 1000);
        return <TimeAgo datetime={lastDate} locale="zh_CN" />;
      },
    },
  ];

  const getRowProps = (item: { id: string }) => {
    const { id } = item;
    return {
      "data-test-subj": `row-${id}`,
      className: "customRowClass",
      // onClick: () => {},
    };
  };

  const getCellProps = (item: { id: string }, column: { field: string }) => {
    const { id } = item;
    const { field } = column;
    return {
      className: "customCellClass",
      "data-test-subj": `cell-${id}-${field}`,
      textOnly: true,
    };
  };
  return (
    <Suspense fallback={<Loading lines={5} />}>
      <EuiPanel hasBorder={true} paddingSize={"s"}>
        <EuiBasicTable
          tableCaption="Git files"
          items={data!}
          rowHeader="name"
          columns={columns}
          rowProps={getRowProps}
          cellProps={getCellProps}
        />
      </EuiPanel>
    </Suspense>
  );
}
