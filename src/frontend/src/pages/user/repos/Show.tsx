import React from "react";
import { useTitle } from "react-use";
import { useParams } from "react-router-dom";

import { RepositoryDetail } from "../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../core/common/document";
import { useGlobal } from "../../../core/global/init";
import i18n from "../../../core/i18n/i18n";
import Error404 from "../../common/404";

export default function RepositoryShow() {
  useTitle(getTitle(i18n.t("repository.menu")));
  const global = useGlobal();
  const currentUser = global.currentUser;

  let namespace = currentUser?.namespace;

  const { repoPath, namespacePath } = useParams();
  useTitle(getTitle(repoPath));

  if (namespacePath !== undefined) {
    namespace = namespacePath;
  }

  if (repoPath === undefined || namespace === undefined) {
    return <Error404 />;
  }

  return (
    <div>
      <RepositoryDetail namespace={namespace} repo={repoPath} />
    </div>
  );
}
