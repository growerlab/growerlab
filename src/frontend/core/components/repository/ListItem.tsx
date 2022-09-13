import React from "react";
import Link from "next/link";

import { repoIcon, repoPath } from "./common";
import { RepositoryEntity } from "../../api/repository/types";
import { Router } from "../../config/router";

interface Args {
  repo: RepositoryEntity
}

export function ListItem(props: Args) {
  const { repo } = props;

  return (
    <div>
      <Link href={repoPath(repo.owner, repo.path)}>
        <a>
          {repoIcon(repo.public)}
          {repo.name}
        </a>
      </Link>
      <div>{repo.description}</div>
    </div>
  );
}
