import React from "react";
import { useTitle } from "react-use";
import { EuiTitle } from "@elastic/eui";

import i18n from "../../../core/i18n/i18n";
import { getTitle } from "../../../core/common/document";
import { NewRepositoryFrom } from "../../../core/components/repository/New";

export default function RepositoryNew(props: React.PropsWithChildren<any>) {
  useTitle(getTitle(i18n.t("repository.create_repository")));

  return (
    <div>
      <EuiTitle size={"s"} className={"mb-5"}>
        <div>{i18n.t<string>("repository.create_repository")}</div>
      </EuiTitle>
      <NewRepositoryFrom />
    </div>
  );
}
