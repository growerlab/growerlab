import { AxiosResponse } from "axios";

import { TypeRepository } from "../../common/types";
import { global } from "../../global/init";
import { API, request } from "../../api/api";

export interface RepositoryRequest {
  namespace: string;
  name: string;
  description: string;
  public: boolean;
}

export class Repository {
  namespace: string;

  constructor(namespace: string) {
    this.namespace = namespace;
  }

  create(req: RepositoryRequest): Promise<AxiosResponse> {
    const url = API.Repositories.Create.render({ namespace: this.namespace });
    return request(global.notice!).post<RepositoryRequest, AxiosResponse>(
      url,
      req
    );
  }

  get(repo: string): Promise<AxiosResponse<TypeRepository>> {
    const url = API.Repositories.Detail.render({
      namespace: this.namespace,
      repo: repo,
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
      created_at: 1222,
      public: true,
      namespace: {
        path: "admin",
        owner: {
          name: "admin",
          namespace: "admin",
        },
      },
      git_http_url: "",
      git_ssh_url: "",
      owner: {
        name: "moli",
        namespace: "admin",
      },
    },
  },
];
