import React, { useEffect, useState } from "react";
import { EuiGlobalToastList } from "@elastic/eui";

import { useNotice } from "../../global/state/useNotice";

interface Toast {
  id: string;
  title: string;
  text: string;
  color?: "primary" | "success" | "warning" | "danger";
}

export default function Notice(props: any) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const { notice } = useNotice();

  useEffect(() => {
    if (!notice) {
      return;
    }
    setToasts(toasts.concat(notice));
  }, [notice]);

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
