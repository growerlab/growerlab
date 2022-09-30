import React from "react";
import { withTranslation } from "react-i18next";

import Activate from "../../core/components/user/Activate";
import { getTitle } from "../../core/common/document";

const activateUser = function (props: any) {
  const { t } = props;
  getTitle(t("website.activate_user"));

  return (
    <div>
      <div className="grid grid-cols-6 gap-4">
        <div className="col-start-2 col-span-4 ">
          <Activate />
        </div>
      </div>
    </div>
  );
};

export default withTranslation()(activateUser);
