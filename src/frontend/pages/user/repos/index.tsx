import React, { useEffect, useState } from "react";

import { RepositoryList } from "../../../core/components/repository/List";
import { Session } from "../../../core/services/auth/session";
import i18n from "../../../core/i18n/i18n";
import UserLayout from "../../layouts/user";

export default function Index() {
  const [namespace, setNamespace] = useState<string>("");

  useEffect((): void => {
    Session.getUserInfo().then((info) => {
      setNamespace(info.namespace_path);
    });
  }, []);

  return (
    <div>
      <UserLayout title={i18n.t("repository.menu")}>
        <RepositoryList ownerPath={namespace} />
      </UserLayout>
    </div>
  );
}
