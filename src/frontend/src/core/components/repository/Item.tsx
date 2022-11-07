import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { EuiHorizontalRule } from "@elastic/eui";

import { repoIcon, repoPath, builtInRepoPath } from "./common";
import { RepositoryEntity } from "../../services/repository/types";
import { globalTypes } from "../../global/init";

interface Args {
  global: globalTypes;
  repo: RepositoryEntity;
}

export function Item(props: Args) {
  const { global, repo } = props;

  const [path, setPath] = useState(repoPath(repo.owner.namespace, repo.path));

  useEffect(() => {
    const currentUser = global.currentUser;
    if (currentUser !== undefined) {
      if (currentUser.namespace_path === repo.owner.namespace) {
        setPath(builtInRepoPath(repo.path));
      }
    }
  }, []);

  return (
    <div>
      <Link to={path}>
        <a className={"text-xl font-bold"}>
          {repoIcon(repo.public)}
          <span className={"ml-3"}>{repo.name}</span>
        </a>
      </Link>
      <div className={"text-slate-500 mt-5"}>{repo.description}</div>
      <EuiHorizontalRule />
    </div>
  );
}
