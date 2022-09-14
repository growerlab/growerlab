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
      name: "1",
      path: "repo1",
      description: "好了浪迹顺利打开降落伞肯德基镂空设计旅客大幅减少了咖啡就算了就；阿拉山口多久；阿拉山口就",
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
}
