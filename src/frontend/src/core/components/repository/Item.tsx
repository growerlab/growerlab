import React from "react";
import { Link } from "react-router-dom";

import { repoIcon, publicRepoPath, builtInRepoPath } from "./common";
import { RepositoryEntity } from "../../common/types";
import { GlobalObject } from "../../global/global";

interface Args {
  global: GlobalObject;
  repo: RepositoryEntity;
}

export function Item(props: Args) {
  const { global, repo } = props;

  let path = publicRepoPath(repo.owner.username, repo.path);

  const currentUser = global.currentUser;
  if (
    currentUser !== undefined &&
    currentUser.namespace === repo.owner.username
  ) {
    path = builtInRepoPath(repo.path);
  }

  return (
    <div>
      <Link to={path} className={"text-xl font-bold"}>
        {repoIcon(repo.public)}
        <span className={"ml-3"}>{repo.namespace.path + "/" + repo.name}</span>
      </Link>
      <div className={"text-slate-500 mt-2 ml-12 text-sm"}>
        {repo.description}
      </div>
    </div>
  );
}
