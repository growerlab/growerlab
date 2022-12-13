import React from "react";
import { useTitle } from "react-use";

import Activate from "../../core/components/user/Activate";
import { getTitle } from "../../core/common/document";
import i18n from "../../core/i18n/i18n";
import Notfound404 from "../../core/components/common/404";
import { useParams } from "react-router-dom";

export default function ActivateUser(props: any) {
  useTitle(getTitle(i18n.t("website.activate_user")));
  const { code } = useParams();
  if (code === undefined) {
    return <Notfound404 />;
  }

  return (
    <div>
      <div className="grid grid-cols-6 gap-4">
        <div className="col-start-2 col-span-4 ">
          <Activate code={code} />
        </div>
      </div>
    </div>
  );
}
