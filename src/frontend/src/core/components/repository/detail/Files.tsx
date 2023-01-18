import React, { useState } from "react";
import useSWRImmutable, { Fetcher } from "swr";
import { EuiBasicTable, EuiIcon, EuiLink, EuiPanel } from "@elastic/eui";
import { useTitle } from "react-use";
import TimeAgo from "timeago-react";

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
import { is } from "@elastic/eui/src/utils/prop_types/is";

interface Props extends RepositoryPathGroup {
  reference: string;
  filePath: string;
  repository?: RepositoryEntity;
  onChangeFilePath: (filePath: string) => void;
  onChangeReference: (reference: string) => void;
}

export function Files(props: Props) {
  const {
    namespace,
    repo,
    reference,
    filePath,
    repository,
    onChangeFilePath,
    onChangeReference,
  } = props;

  useTitle(getTitle(repo));

  const repositoryAPI = useRepositoryAPI(namespace);
  const [fileEntities, setFileEntities] = useState<FileEntity[]>();
  const [isEmptyTree, setTreeEmpty] = useState<boolean>(false);
  const [repoFilePath, setRepoFilePath] = useState(filePath); // 正在访问的repo路径

  if (!isEmptyTree && repository?.last_push_at == 0) {
    setTreeEmpty(true);
  }

  const fetcher: Fetcher = () => {
    const params = { namespace, repo, ref: reference, dir: filePath };
    return repositoryAPI.treeFiles(params).then((res) => {
      setFileEntities(res.data);
      return res.data;
    });
  };
  useSWRImmutable(
    isEmptyTree ? null : `/swr/key/repo/${namespace}/${repo}/tree_files`,
    fetcher,
    { shouldRetryOnError: false }
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
  if (!fileEntities) {
    return <Loading lines={5} />;
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
        return (
          <EuiLink href="#文件内容">
            {icon} {name}
          </EuiLink>
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
    <>
      <EuiPanel hasBorder={true} paddingSize={"s"}>
        <EuiBasicTable
          tableCaption="Git files"
          items={fileEntities}
          rowHeader="name"
          columns={columns}
          rowProps={getRowProps}
          cellProps={getCellProps}
        />
      </EuiPanel>
    </>
  );
}
