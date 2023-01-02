import React, { useState } from "react";
import useSWRImmutable, { Fetcher } from "swr";
// import { useSearchParams } from "react-router-dom";
import { formatDate, EuiBasicTable, EuiLink, EuiPanel } from "@elastic/eui";
import { useSearchParams } from "react-router-dom";
import { useTitle } from "react-use";

import {
  FileEntity,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import Loading from "../../common/Loading";
import { useRepositoryAPI } from "../../../api/repository";
import EmptyTree from "./EmptyTree";
import { getTitle } from "../../../common/document";

interface Props extends RepositoryPathGroup {
  reference: string;
  filePath: string;
  repository?: RepositoryEntity;
  onChangeFilePath: (filePath: string) => void;
}

export function Files(props: Props) {
  const { namespace, repo, reference, filePath, repository, onChangeFilePath } =
    props;

  useTitle(getTitle(repo));

  const repositoryAPI = useRepositoryAPI(namespace);
  const [fileEntities, setFileEntities] = useState<FileEntity[]>();
  const [isEmptyTree, setTreeEmpty] = useState<boolean>(false);
  const [repoFilePath, setRepoFilePath] = useState(filePath);

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
      name: "Name",
      sortable: true,
      render: (name: string) => <EuiLink href="#">{name}</EuiLink>,
    },
    {
      field: "lastCommitMessage",
      name: "Last commit message",
      render: (item: any) => <EuiLink href="#">{item}</EuiLink>,
    },
    {
      field: "lastCommitDate",
      name: "Last commit date",
      dataType: "date",
      render: (date: Date) => formatDate(date, "dobLong"),
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
