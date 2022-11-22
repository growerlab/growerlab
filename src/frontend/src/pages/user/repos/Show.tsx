import React from "react";
import { useTitle } from "react-use";
import { useOutlet } from "react-router-dom";

import { getTitle } from "../../../core/common/document";
import i18n from "../../../core/i18n/i18n";
import Error404 from "../../common/404";
import { useRepositoryPathGroup } from "../../../core/components/hook/repository";

interface Props extends React.PropsWithChildren {
  defaultChild?: React.ReactElement;
}

export default function RepositoryShow(props: Props) {
  useTitle(getTitle(i18n.t("repository.menu")));

  const outlet = useOutlet();
  const defaultOutlet = outlet === null ? props.defaultChild : outlet;

  const { isInvalid } = useRepositoryPathGroup();
  if (isInvalid) {
    return <Error404 />;
  }

  return <div>{defaultOutlet}</div>;
}
