import React, { useEffect, useState } from "react";

import { RepositoryList } from "../../../components/repository/List";
import { Session, LoginInfo } from "../../../services/auth/session";
import i18n from "../../../i18n/i18n";
import UserLayout from "../../layouts/user";

export default function Index() {
  const [namespace, setNamespace] = useState<string>("");

  useEffect((): void => {
    Session.getUserInfo().then((info) => {
      setNamespace(info.namespacePath);
    });
  }, []);

  return (
    <div>
      <UserLayout title={i18n.t("repository.list")}>
        <RepositoryList ownerPath={namespace} />
      </UserLayout>
    </div>
  );
}
