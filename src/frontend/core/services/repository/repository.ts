import { AxiosResponse } from "axios";

import { TypeRepository } from "./types";
import { global } from "../../global/init";
import { API, request } from "../../api/api";

export class Repository {
  ownerPath: string

  constructor(ownerPath: string) {
    this.ownerPath = ownerPath;
  }

  get(repoPath: string): Promise<AxiosResponse<TypeRepository>> {
    return request(global.notice!)
      .get<TypeRepository, AxiosResponse<TypeRepository>>(API.Repositories.Detail.render({ owenrPath: this.ownerPath, repoPath: repoPath }))
      .then((res) => {
        return res;
      });
  }

  list(page = 0): TypeRepository[] | undefined {
    return mockRepositories;
  }
}

const mockRepositories: TypeRepository[] = [
  {
    repository: {
      uuid: "1",
      name: "hello",
      path: "repo1",
      description:
        "这是一个仓库描述这是一个仓库描述这是一个仓库描述这是一个仓库描述;这是一个仓库描述；这是一个仓库描述",
      createdAt: 1222,
      public: true,
      // pathGroup: "123123/sdfsdf",
      gitHttpURL: "",
      gitSshURL: "",
      owner: {
        name: "moli",
        namespace: "admin",
      },
    },
  }
];
