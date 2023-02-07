import React from "react";
import useSWR from "swr";
import TimeAgo from "timeago-react";
import { EuiBasicTable, EuiIcon, EuiLink } from "@elastic/eui";

import {
  DetailType,
  FileEntity,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import { useRepositoryAPI } from "../../../api/repository";
import i18n from "../../../i18n/i18n";
import { Router } from "../../../../config/router";
import { Path } from "../../../common/path";

interface Props extends RepositoryPathGroup {
  reference: string;
  initialFolder: string;
  repository?: RepositoryEntity;
  onChange: (name: string, type: DetailType) => void;
}

export default function Tree(props: Props) {
  const { namespace, repo, reference, initialFolder, repository } = props;

  const repositoryAPI = useRepositoryAPI(namespace);
  const currentRepoFolder = new Path(initialFolder); // 正在访问的repo路径

  const buildTreeURL = (tree: string) => {
    return Router.User.Repository.Reference.render({
      refType: "tree",
      "*": tree,
      ref: reference,
      repo: repo,
    });
  };
  const buildBlobURL = (filePath: string) => {
    return Router.User.Repository.Reference.render({
      refType: "blob",
      "*": filePath,
      ref: reference,
      repo: repo,
    });
  };

  const icons = {
    document: <EuiIcon type={"document"} />,
    folderClosed: <EuiIcon type={"folderClosed"} />,
  };

  // tree api
  const fetcherKey = `/swr/key/repo/${namespace}/${repo}/${window.location.pathname}`;
  const treeFetcher = () => {
    const params = {
      namespace,
      repo,
      ref: reference,
      folder: currentRepoFolder.toString(),
    };
    return repositoryAPI.treeFiles(params).then((res) => {
      return res.data;
    });
  };
  const { data } = useSWR<FileEntity[]>(fetcherKey, treeFetcher, {
    suspense: true,
  });

  const columns: any = [
    {
      field: "name",
      name: i18n.t("repository.commit.file_name"),
      render: (name: string, record: FileEntity) => {
        const icon = record.is_file ? icons.document : icons.folderClosed;
        const type: DetailType = record.is_file ? "blob" : "tree";
        const folderPath = currentRepoFolder.toString() + "/" + name;
        const link = record.is_file
          ? buildBlobURL(folderPath)
          : buildTreeURL(folderPath);
        return (
          <>
            {icon}{" "}
            <EuiLink
              href={link}
              onClick={() => {
                // navigate(link, { state: link });
                // currentRepoFolder.append(name);
                props.onChange(name, type);
              }}
            >
              {name}
            </EuiLink>
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
    };
  };

  const getCellProps = (item: { id: string }, column: { field: string }) => {
    const { id } = item;
    const { field } = column;
    return {
      "data-test-subj": `cell-${id}-${field}`,
      textOnly: true,
    };
  };

  return (
    <>
      <EuiBasicTable
        tableCaption="Git files"
        items={data!}
        rowHeader="name"
        columns={columns}
        rowProps={getRowProps}
        cellProps={getCellProps}
      />
    </>
  );
}
