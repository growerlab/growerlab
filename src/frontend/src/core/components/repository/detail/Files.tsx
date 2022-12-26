import React, { useState } from "react";
import useSWR, { Fetcher } from "swr";

import { FileEntity, RepositoryPathGroup } from "../../../common/types";
import Loading from "../../common/Loading";
import { useRepositoryAPI } from "../../../api/repository";

interface Props extends RepositoryPathGroup {
  reference: "master" | string;
  dir: "" | string;
}

export function Files(props: Props) {
  const { namespace, repo, reference, dir } = props;
  const repositoryAPI = useRepositoryAPI(namespace);
  const [fileEntities, setFileEntities] = useState<FileEntity[]>();

  const fetcher: Fetcher = () => {
    const params = { namespace, repo, ref: reference, dir };
    return repositoryAPI.treeFiles(params).then((res) => {
      setFileEntities(res.data);
      return res.data;
    });
  };
  useSWR(`/swr/key/repo/${namespace}/${repo}/tree_files`, fetcher);
  if (!fileEntities) {
    return <Loading lines={5} />;
  }

  return (
    <>
      <Loading lines={5}></Loading>
    </>
  );
}
