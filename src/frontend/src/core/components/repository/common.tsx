import React from "react";
import { EuiAvatar } from "@elastic/eui";
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

export function builtInRepoPath(repoPath: string): string {
  return Router.User.Repository.Show.render({
    repo: repoPath,
  });
}

// 公共路径
export function repoPath(namespace: string, repoPath: string): string {
  return Router.Namespace.Repository.render({
    namespace: namespace,
    repo: repoPath,
  });
}
