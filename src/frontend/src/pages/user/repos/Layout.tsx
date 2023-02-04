import React, { useEffect } from "react";
import { useOutlet } from "react-router-dom";
import { useUserMenu } from "../../../core/global/state/useMenu";

interface Props extends React.PropsWithChildren {
  defaultChild?: React.ReactElement;
}

export default function RepositoryLayout(props: Props) {
  const outlet = useOutlet();
  const defaultOutlet = outlet === null ? props.defaultChild : outlet;

  // 设置user下的menu后台选项
  const { current, setMenuSelect } = useUserMenu();
  useEffect(() => {
    setMenuSelect("repository");
  }, [current]);

  return <>{defaultOutlet}</>;
}
