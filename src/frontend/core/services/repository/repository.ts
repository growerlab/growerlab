import { TypeRepositories, TypeRepository } from "./types";

export class Repository {
  ownerPath: string

  constructor(ownerPath: string) {
    this.ownerPath = ownerPath;
    if (this.ownerPath === undefined) {
      // const current = getUserInfo();
      // if (current !== null) {
      //   this.repo.ownerPath = current.namespacePath;
      // }
    }
  }

  get(repoPath: string): TypeRepository | null {
    console.info(repoPath)
    return null;
  }

  list(page = 0): TypeRepositories | null {
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
        namespace: "admin",
      },
    },
  ],
};
