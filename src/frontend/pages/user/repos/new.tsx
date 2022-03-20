import React from "react";

import { NewRepositoryFrom } from "../../../components/repository/New";
import i18n from "../../../i18n/i18n";
import { getTitle } from "../../../common/document";

export default function (props: React.PropsWithChildren<any>) {
  getTitle(i18n.t("repository.create_repository"));

  return (
    <div>
      <NewRepositoryFrom />
    </div>
  );
}
