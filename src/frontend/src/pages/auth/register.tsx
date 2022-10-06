import React from "react";

import { getTitle } from "../../core/common/document";
import RegisterForm from "../../core/components/user/Register";
import i18n from "../../core/i18n/i18n";

export default function Register(props: any) {
  getTitle(i18n.t("website.register"));

  return (
    <div>
      <div className="grid grid-cols-6 gap-4">
        <div className="col-start-2 col-span-4 ">
          <RegisterForm />
        </div>
      </div>
    </div>
  );
}
