import React from "react";
import { useRouter } from "next/router";
import Head from "next/head";
import { EuiLoadingSpinner } from "@elastic/eui";

import UserLayout from "../../layouts/user";
import { RepositoryDetail } from "../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../core/common/document";
import i18n from "../../../core/i18n/i18n";
import { useGlobal } from "../../../core/global/init";

export default function ShowRepoPage() {
  const global = useGlobal();
  const router = useRouter();
  const { repoPath } = router.query;

  console.info("== = = == ", router.query);
  const user = global.getUserInfo();
  if (user == null) {
    return <EuiLoadingSpinner size="xl" />;
  }

  return (
    <div>
      <Head>
        <title>{getTitle(i18n.t(repoPath))}</title>
      </Head>
      <UserLayout title={i18n.t("repository.menu")}>
        <React.Suspense fallback={<EuiLoadingSpinner size="xl" />}>
          <RepositoryDetail
            currentUser={user}
            ownerPath={user?.namespace_path}
            repoPath={repoPath}
          />
        </React.Suspense>
      </UserLayout>
    </div>
  );
}
