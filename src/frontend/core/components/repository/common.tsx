import React from "react";
import { EuiIcon, EuiAvatar } from "@elastic/eui";
import { Owner } from "../../services/repository/types";
import { Router } from "../../../config/router";

const repoPrivateIcon = (
  <EuiAvatar name="lock" iconType="lock" color="#E6F1FA" />
);
const repoPublicIcon = (
  <EuiAvatar name="lock" iconType="lockOpen" color="#E6F1FA" />
);

export function repoIcon(pub: boolean) {
  return pub ? repoPublicIcon : repoPrivateIcon;
}

export function repoPath(owner: Owner, path: string): string {
  return Router.Namespace.Repository.render({
    namespacePath: owner.namespace,
    repoPath: path,
  });
}
