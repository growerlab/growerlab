
import { LockIcon, UnlockIcon } from "evergreen-ui";

import { Router } from "../../config/router";
import { Owner } from "../../api/repository/types";

const repoPrivateIcon = <LockIcon />;
const repoPublicIcon = <UnlockIcon />;

export function repoIcon(pub: boolean) {
  return pub ? repoPublicIcon : repoPrivateIcon;
}

export function repoPath(owner: Owner, path: string): string {
  return Router.Namespace.Repository.render({ namespacePath: owner.namespace, repoPath: path });
}