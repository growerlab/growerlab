import React from "react";

import i18n from "../../../core/i18n/i18n";
import { NewRepositoryForm } from "../../../core/components/repository/NewRepositoryForm";
import Header from "../../../core/components/ui/common/Header";
import { useGlobal } from "../../../core/global/global";
import { useTitle } from "../../../core/global/state";

export default function RepositoryNew() {
  useTitle(i18n.t("repository.create_repository"));
  const { currentUser } = useGlobal();
  if (currentUser === undefined) return <></>;

  return (
    <div className="p-5">
      <Header title={i18n.t<string>("repository.create_repository")} />
      <NewRepositoryForm namespace={currentUser.namespace} />
    </div>
  );
}
