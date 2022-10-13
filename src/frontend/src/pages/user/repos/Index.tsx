import React, { useEffect, useState } from "react";
import { useTitle } from "react-use";

import { RepositoryList } from "../../../core/components/repository/List";
import { Session } from "../../../core/services/auth/session";
import i18n from "../../../core/i18n/i18n";

export default function RepositoryIndex() {
  useTitle(i18n.t("repository.menu"));
  const user = Session.getUserInfo();

  return (
    <>
      <RepositoryList ownerPath={user.namespace_path} />
    </>
  );
}
