import React, { useEffect, useState } from "react";
import { useTitle } from "react-use";

import { RepositoryList } from "../../../core/components/repository/List";
import { Session } from "../../../core/services/auth/session";
import i18n from "../../../core/i18n/i18n";

export default function RepositoryIndex() {
  const [namespace, setNamespace] = useState<string>("");
  useTitle(i18n.t("repository.menu"));

  useEffect((): void => {
    const user = Session.getUserInfo();
    if (user !== undefined) {
      setNamespace(user.namespace_path);
    }
  }, []);

  return (
    <>
      <RepositoryList ownerPath={namespace} />
    </>
  );
}
