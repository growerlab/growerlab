import React from "react";
import { useTitle } from "react-use";

import i18n from "../../../core/i18n/i18n";
import { getTitle } from "../../../core/common/document";
import { NewRepositoryFrom } from "../../../core/components/repository/New";
import Title from "../../../core/components/common/Title";

export default function RepositoryNew(props: React.PropsWithChildren<any>) {
  useTitle(getTitle(i18n.t("repository.create_repository")));

  return (
    <div>
      <Title title={i18n.t<string>("repository.create_repository")} />
      <NewRepositoryFrom />
    </div>
  );
}
