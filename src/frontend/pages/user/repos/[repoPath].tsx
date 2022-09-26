import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";
import Head from "next/head";

import UserLayout from "../../layouts/user";
import { RepositoryDetail } from "../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../core/common/document";
import i18n from "../../../core/i18n/i18n";
import { useGlobal } from "../../../core/global/init";
import { LoginInfo } from "../../../core/services/auth/session";

export default function ShowRepoPage() {
  const global = useGlobal();
  const router = useRouter();
  const repoPath = router.query.repoPath as string;
  const [ownerPath, setOwnerPath] = useState<string>("");
  const [currentUser, setCurrentUser] = useState<LoginInfo>();

  useEffect(() => {
    global.getUserInfo().then((user) => {
      setCurrentUser(user);
      setOwnerPath(user.namespace_path);
    });
  }, []);

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
