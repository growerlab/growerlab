import React from "react";

import { RepositoryDetail } from "../../../../core/components/repository/RepositoryDetail";
import i18n from "../../../../core/i18n/i18n";
import { useRepositoryPathGroup } from "../../../../core/components/hook/repository";
import { useTitle } from "../../../../core/global/state";

type Props = any;

export default function Detail(props: Props) {
  useTitle(i18n.t("repository.menu"));

  const { namespace, repo } = useRepositoryPathGroup();

  return (
    <div>
      <RepositoryDetail namespace={namespace} repo={repo} {...props} />
    </div>
  );
}
