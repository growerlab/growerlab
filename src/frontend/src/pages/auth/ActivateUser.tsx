import React from "react";

import Activate from "../../core/components/user/Activate";
import { getTitle } from "../../core/common/document";
import i18n from "../../core/i18n/i18n";

export default function ActivateUser(props: any) {
  getTitle(i18n.t("website.activate_user"));

  return (
    <div>
      <div className="grid grid-cols-6 gap-4">
        <div className="col-start-2 col-span-4 ">
          <Activate />
        </div>
      </div>
    </div>
  );
}
