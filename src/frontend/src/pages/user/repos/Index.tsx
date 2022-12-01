import React from "react";
import { useTitle } from "react-use";

import { RepositoryList } from "../../../core/components/repository/List";
import i18n from "../../../core/i18n/i18n";
import { useGlobal } from "../../../core/global/global";

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
