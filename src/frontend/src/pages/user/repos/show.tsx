import React from "react";
import { EuiLoadingSpinner } from "@elastic/eui";
import { useParams } from "react-router-dom";

import UserLayout from "../../layouts/UserLayout";
import { RepositoryDetail } from "../../../core/components/repository/RepositoryDetail";
import { getTitle } from "../../../core/common/document";
import i18n from "../../../core/i18n/i18n";
import { useGlobal } from "../../../core/global/init";
import Error404 from "../../common/404";

export default function ShowRepoPage() {
  const global = useGlobal();
  const { repoPath } = useParams();
  if (repoPath === undefined) {
    return <Error404 />;
  }

  const user = global.getUserInfo();
  if (user == null) {
    return <EuiLoadingSpinner size="xl" />;
  }

  return (
    <div>
      <div>{getTitle(i18n.t(repoPath))}</div>
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
