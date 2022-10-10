import React from "react";
import { useTitle } from "react-use";

import i18n from "../../../core/i18n/i18n";
import { getTitle } from "../../../core/common/document";
import { NewRepositoryFrom } from "../../../core/components/repository/New";

export default function New(props: React.PropsWithChildren<any>) {
  useTitle(getTitle(i18n.t("repository.create_repository")));

  return (
    <div>
      <NewRepositoryFrom />
    </div>
  );
}
