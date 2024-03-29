import React from "react";

import RegisterForm from "../../core/components/user/RegisterForm";
import i18n from "../../core/i18n/i18n";
import { useTitle } from "../../core/global/state";

export default function Register(props: any) {
  useTitle(i18n.t("website.register"));

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
