import React from "react";
import { useTitle } from "react-use";

import i18n from "../../../core/i18n/i18n";
import { getTitle } from "../../../core/common/document";
import { NewRepositoryFrom } from "../../../core/components/repository/NewRepositoryFrom";
import Title from "../../../core/components/common/Title";
import { Session } from "../../../core/services/auth/session";

export default function RepositoryNew(props: React.PropsWithChildren<any>) {
  useTitle(getTitle(i18n.t("repository.create_repository")));
  const user = Session.getUserInfo();

  return (
    <div>
      <Title title={i18n.t<string>("repository.create_repository")} />
      <NewRepositoryFrom ownerPath={user.namespace_path} />
    </div>
  );
}
