import React from "react";
import { withTranslation } from "react-i18next";

import { getTitle } from "../../core/common/document";
import RegisterForm from "../../core/components/user/Register";

const register = function Register(props: any) {
  const { t } = props;
  getTitle(t("website.register"));

  return (
    <div>
      <div className="grid grid-cols-6 gap-4">
        <div className="col-start-2 col-span-4 ">
          <RegisterForm />
        </div>
      </div>
    </div>
  );
};

export default withTranslation()(register);
