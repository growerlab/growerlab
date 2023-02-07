import React from "react";

import Activate from "../../core/components/user/Activate";
import i18n from "../../core/i18n/i18n";
import Notfound404 from "../../core/components/ui/common/404";
import { useParams } from "react-router-dom";
import { useTitle } from "../../core/global/state";

export default function ActivateUser(props: any) {
  useTitle(i18n.t("website.activate_user"));
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
