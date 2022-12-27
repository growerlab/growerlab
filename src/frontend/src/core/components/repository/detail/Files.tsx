import React, { useState } from "react";
import useSWRImmutable, { Fetcher } from "swr";

import { formatDate, EuiBasicTable, EuiLink } from "@elastic/eui";

import { FileEntity, RepositoryPathGroup } from "../../../common/types";
import Loading from "../../common/Loading";
import { useRepositoryAPI } from "../../../api/repository";
import EmptyTree from "./EmptyTree";

interface Props extends RepositoryPathGroup {
  reference: "master" | string;
  dir: "" | string;
  isEmptyTree: boolean;
}

export function Files(props: Props) {
  const { namespace, repo, reference, dir, isEmptyTree } = props;
  const repositoryAPI = useRepositoryAPI(namespace);
  const [fileEntities, setFileEntities] = useState<FileEntity[]>();

  const fetcher: Fetcher = () => {
    const params = { namespace, repo, ref: reference, dir };
    return repositoryAPI.treeFiles(params).then((res) => {
      setFileEntities(res.data);
      return res.data;
    });
  };
  const { error } = useSWRImmutable(
    isEmptyTree ? null : `/swr/key/repo/${namespace}/${repo}/tree_files`,
    fetcher,
    { shouldRetryOnError: false }
  );
  if (isEmptyTree) {
    return (
      <EmptyTree
        cloneURLSSH={"ssh://git.com"}
        cloneURLHttp={"https://git.com"}
        defaultBranch={"main"}
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
