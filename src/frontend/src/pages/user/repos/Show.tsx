import React from "react";
import { useTitle } from "react-use";
import { useParams } from "react-router-dom";

import { RepositoryDetail } from "../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../core/common/document";
import { useGlobal } from "../../../core/global/init";
import i18n from "../../../core/i18n/i18n";
import Error404 from "../../common/404";
import { useRepositoryPathGroup } from "../../../core/components/hook/repository";

export default function RepositoryShow() {
  useTitle(getTitle(i18n.t("repository.menu")));

  const { namespace, repo, isInvalid } = useRepositoryPathGroup();

  if (isInvalid) {
    return <Error404 />;
  }

  return (
    <div>
      <RepositoryDetail namespace={namespace} repo={repo} />
    </div>
  );
}
