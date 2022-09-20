import { RepositoryArgs, TypeRepositories, TypeRepository } from "./types";

export class Repository {
  repo: RepositoryArgs;

  constructor(args: RepositoryArgs) {
    this.repo = args;
    if (this.repo.ownerPath === undefined) {
      // const current = getUserInfo();
      // if (current !== null) {
      //   this.repo.ownerPath = current.namespacePath;
      // }
    }
  }

  get(repoPath: string): TypeRepository | null {
    return null;
  }

  list(): TypeRepositories | null {
    return mockRepositories;
  }
}

const mockRepositories: TypeRepositories = {
  repositories: [
    {
      uuid: "1",
      name: "hello",
      path: "repo1",
      description:
        "这是一个仓库描述这是一个仓库描述这是一个仓库描述这是一个仓库描述;这是一个仓库描述；这是一个仓库描述",
      createdAt: 1222,
      public: true,
      pathGroup: "",
      gitHttpURL: "",
      gitSshURL: "",
      owner: {
        name: "moli",
        namespace: "moli",
      },
    },
  ],
};
