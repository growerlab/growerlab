import React, { useState } from "react";
import useSWRImmutable, { Fetcher } from "swr";

import { formatDate, EuiBasicTable, EuiLink } from "@elastic/eui";

import {
  FileEntity,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import Loading from "../../common/Loading";
import { useRepositoryAPI } from "../../../api/repository";
import EmptyTree from "./EmptyTree";

interface Props extends RepositoryPathGroup {
  reference: string;
  dir: string;
  repository?: RepositoryEntity;
}

export function Files(props: Props) {
  const { namespace, repo, reference, dir, repository } = props;
  const repositoryAPI = useRepositoryAPI(namespace);
  const [fileEntities, setFileEntities] = useState<FileEntity[]>();
  const [isEmptyTree, setTreeEmpty] = useState<boolean>(false);

  if (!isEmptyTree && repository?.last_push_at == 0) {
    setTreeEmpty(true);
  }

  const fetcher: Fetcher = () => {
    const params = { namespace, repo, ref: reference, dir };
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
      field: "firstName",
      name: "First Name",
      sortable: true,
      "data-test-subj": "firstNameCell",
      mobileOptions: {
        render: (item: any) => (
          <span>
            {item.firstName}{" "}
            <EuiLink href="#" target="_blank">
              {item.lastName}
            </EuiLink>
          </span>
        ),
        header: false,
        truncateText: false,
        enlarge: true,
        width: "100%",
      },
    },
    {
      field: "lastName",
      name: "Last Name",
      truncateText: true,
      render: (name: string) => (
        <EuiLink href="#" target="_blank">
          {name}
        </EuiLink>
      ),
      mobileOptions: {
        show: false,
      },
    },
    {
      field: "dateOfBirth",
      name: "Date of Birth",
      dataType: "date",
      render: (date: Date) => formatDate(date, "dobLong"),
    },
  ];

  const items: any = []; //store.users.filter((user, index) => index < 10);

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
      <EuiBasicTable
        tableCaption="Demo of EuiBasicTable"
        items={items}
        rowHeader="firstName"
        columns={columns}
        rowProps={getRowProps}
        cellProps={getCellProps}
      />
    </>
  );
}
