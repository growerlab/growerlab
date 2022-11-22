import React from "react";
import { useTitle } from "react-use";

import { RepositoryDetail } from "../../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../../core/common/document";
import i18n from "../../../../core/i18n/i18n";
import { useRepositoryPathGroup } from "../../../../core/components/hook/repository";

export default function Detail(props: any) {
  useTitle(getTitle(i18n.t("repository.menu")));

  const { namespace, repo } = useRepositoryPathGroup();

  return (
    <div>
      <RepositoryDetail namespace={namespace} repo={repo} />
    </div>
  );
}
