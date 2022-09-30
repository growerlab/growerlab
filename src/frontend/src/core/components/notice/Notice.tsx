import React, { useEffect, useState } from "react";
import { EuiGlobalToastList } from "@elastic/eui";
import { withTranslation } from "react-i18next";
import { useRecoilValue } from "recoil";

import { noticeState } from "../../global/recoil/notice";

interface Toast {
  id: string;
  title: string;
  text: string;
  color?: "primary" | "success" | "warning" | "danger";
}

function Notice(props: any) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const value = useRecoilValue(noticeState);

  useEffect(
    function () {
      if (value === undefined) {
        return;
      }
      setToasts(toasts.concat(value));
    },
    [value]
  );

  return (
    <EuiGlobalToastList
      toasts={toasts}
      dismissToast={(removedToast) => {
        setToasts(toasts.filter((toast) => toast.id !== removedToast.id));
      }}
      toastLifeTimeMs={3000}
    />
  );
}

export default withTranslation()(Notice);
