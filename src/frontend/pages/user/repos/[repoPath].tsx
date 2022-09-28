import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";
import Head from "next/head";

import UserLayout from "../../layouts/user";
import { RepositoryDetail } from "../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../core/common/document";
import i18n from "../../../core/i18n/i18n";
import { useGlobal } from "../../../core/global/init";
import { UserInfo } from "../../../core/services/auth/session";

export default function ShowRepoPage() {
  const global = useGlobal();
  const router = useRouter();
  const repoPath = router.query.repoPath as string;
  const [ownerPath, setOwnerPath] = useState<string>("");
  const [currentUser, setCurrentUser] = useState<UserInfo>();

  useEffect(() => {
    const user = global.getUserInfo();
    if (user !== undefined) {
      setCurrentUser(user);
      setOwnerPath(user.namespace_path);
    } else {
      return undefined;
    }
  }, []);

  if (
    currentUser === undefined ||
    ownerPath.length == 0 ||
    repoPath.length == 0
  ) {
    return <h1>404</h1>;
  }

  return (
    <div>
      <Head>
        <title>{getTitle(i18n.t(repoPath))}</title>
      </Head>
      <UserLayout title={i18n.t("repository.menu")}>
        <RepositoryDetail
          currentUser={currentUser}
          ownerPath={ownerPath}
          repoPath={repoPath}
        />
      </UserLayout>
    </div>
  );
}
