import React, { useState, Suspense } from "react";
import useSWR from "swr";
import { useTitle } from "react-use";
import TimeAgo from "timeago-react";
import { useNavigate } from "react-router-dom";
import { EuiBasicTable, EuiIcon, EuiLink, EuiPanel } from "@elastic/eui";
import { EuiBreadcrumbProps } from "@elastic/eui/src/components/breadcrumbs/breadcrumb";

import {
  FileEntity,
  RepositoryEntity,
  RepositoryPathGroup,
} from "../../../common/types";
import { useRepositoryAPI } from "../../../api/repository";
import EmptyTree from "./EmptyTree";
import { getTitle } from "../../../common/document";
import i18n from "../../../i18n/i18n";
import { Router } from "../../../../config/router";
import Loading from "../../ui/common/Loading";
import Breadcrumbs from "../../ui/Breadcrumbs";
import { Path } from "../../../common/path";

interface Props extends RepositoryPathGroup {
  reference: string;
  initialFolder: string;
  repository?: RepositoryEntity;
}

export function Files(props: Props) {
  const { namespace, repo, reference, initialFolder, repository } = props;

  useTitle(getTitle(repo));

  const navigate = useNavigate();
  const repositoryAPI = useRepositoryAPI(namespace);
  const [isEmptyTree, setTreeEmpty] = useState<boolean>(false);
  const [currentRepoFolder] = useState(new Path(initialFolder)); // 正在访问的repo路径
  if (!isEmptyTree && repository?.last_push_at == 0) {
    setTreeEmpty(true);
  }

  const buildTreeURL = (tree: string) => {
    return Router.User.Repository.Tree.render({
      "*": tree,
      ref: reference,
      repo: repo,
    });
  };
  const buildBlobURL = (filePath: string) => {
    return Router.User.Repository.Blob.render({
      filepath: filePath,
      ref: reference,
      repo: repo,
    });
  };

  const icons = {
    document: <EuiIcon type={"document"} />,
    folderClosed: <EuiIcon type={"folderClosed"} />,
  };

  // api
  const fetcher = () => {
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
        const icon = record.is_file ? icons.document : icons.folderClosed;
        const folderPath = currentRepoFolder.toString() + "/" + name;
        const link = record.is_file
          ? buildBlobURL(name)
          : buildTreeURL(folderPath);
        return (
          <>
            {icon}{" "}
            <EuiLink
              // to={link}
              onClick={() => {
                navigate(link, { state: link });
                currentRepoFolder.append(name);
                return;
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

  // folder paths
  const folderBreadcrumbs: EuiBreadcrumbProps[] = [
    {
      text: "/",
      onClick: () => {
        const rootURL = buildTreeURL("");
        navigate(rootURL, { state: rootURL });
        currentRepoFolder.reset("");
      },
    },
  ];
  currentRepoFolder.forEach((value, index, array) => {
    if (value === "") {
      return;
    }
    const link = buildTreeURL(value);
    const onClick = () => {
      navigate(link, { state: link });
      currentRepoFolder.reset(array.slice(0, index + 1));
    };
    folderBreadcrumbs.push({
      text: value,
      onClick: index == array.length - 1 ? undefined : onClick,
    });
  });

  return (
    <Suspense fallback={<Loading lines={5} />}>
      <EuiPanel hasBorder={true} paddingSize={"s"}>
        <Breadcrumbs
          truncate={true}
          max={4}
          breadcrumbs={folderBreadcrumbs}
          icon={"submodule"}
        />
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
