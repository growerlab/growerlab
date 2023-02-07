import React from "react";

import { RepositoryList } from "../../../core/components/repository/List";
import i18n from "../../../core/i18n/i18n";
import { useGlobal } from "../../../core/global/global";
import { useTitle } from "../../../core/global/state";

export default function RepositoryIndex() {
  useTitle(i18n.t("repository.menu"));
  const { currentUser } = useGlobal();

  if (currentUser === undefined) return <></>;

  return (
    <div>
      <RepositoryList namespace={currentUser.namespace} />
    </div>
  );
}
