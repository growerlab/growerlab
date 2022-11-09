import { AxiosResponse } from "axios";

import { TypeRepository } from "./types";
import { global } from "../../global/init";
import { API, request } from "../../api/api";

export interface RepositoryRequest {
  namespace_path: string;
  name: string;
  description: string;
  public: boolean;
}

export class Repository {
  ownerPath: string;

  constructor(ownerPath: string) {
    this.ownerPath = ownerPath;
  }

  create(req: RepositoryRequest): Promise<AxiosResponse> {
    const url = API.Repositories.Create.render({ ownerPath: this.ownerPath });
    return request(global.notice!).post<RepositoryRequest, AxiosResponse>(
      url,
      req
    );
  }

  get(repoPath: string): Promise<AxiosResponse<TypeRepository>> {
    const url = API.Repositories.Detail.render({
      ownerPath: this.ownerPath,
      repoPath: repoPath,
    });
    return request(global.notice!).get<
      TypeRepository,
      AxiosResponse<TypeRepository>
    >(url);
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
  },
];
