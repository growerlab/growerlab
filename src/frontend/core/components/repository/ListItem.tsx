import React from "react";
import Link from "next/link";
import { EuiHorizontalRule } from "@elastic/eui";

import { repoIcon, repoPath } from "./common";
import { RepositoryEntity } from "../../services/repository/types";

interface Args {
  repo: RepositoryEntity;
}

export function ListItem(props: Args) {
  const { repo } = props;

  return (
    <div>
      <Link href={repoPath(repo.owner, repo.path)}>
        <div className={"text-xl font-bold"}>
          {repoIcon(repo.public)}
          <a className={"ml-3"}>{repo.name}</a>
        </div>
      </Link>
      <div className={"text-slate-500 mt-5"}>{repo.description}</div>
      <EuiHorizontalRule />
    </div>
  );
}
