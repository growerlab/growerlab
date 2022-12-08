import React from "react";
import { useTitle } from "react-use";

import i18n from "../../../core/i18n/i18n";
import { getTitle } from "../../../core/common/document";
import { NewRepositoryForm } from "../../../core/components/repository/NewRepositoryForm";
import Title from "../../../core/components/common/Title";
import { useGlobal } from "../../../core/global/global";

export default function RepositoryNew() {
  useTitle(getTitle(i18n.t("repository.create_repository")));
  const { currentUser } = useGlobal();
  if (currentUser === undefined) return <></>;

  return (
    <div className="p-5">
      <Title title={i18n.t<string>("repository.create_repository")} />
      <NewRepositoryForm namespace={currentUser.namespace} />
    </div>
  );
}
