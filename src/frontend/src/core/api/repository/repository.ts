import { AxiosResponse } from "axios";

import { RepositoryEntity, TypeRepositories } from "../../common/types";
import { API, request } from "../api";
import { GlobalObject, useGlobal } from "../../global/global";

export interface RepositoryRequest {
  namespace: string;
  name: string;
  description: string;
  public: boolean;
}

export function useRepositoryAPI(namespace: string) {
  const global = useGlobal();
  const repo = new Repository(namespace, global);
  return repo;
}

class Repository {
  private readonly namespace: string;
  private readonly global: GlobalObject;

  constructor(namespace: string, global: GlobalObject) {
    this.namespace = namespace;
    this.global = global;
  }

  create(req: RepositoryRequest): Promise<AxiosResponse> {
    const url = API.Repositories.Create.render({ namespace: this.namespace });
    return request(this.global).post<RepositoryRequest, AxiosResponse>(
      url,
      req
    );
  }

  get(repo: string): Promise<AxiosResponse<RepositoryEntity>> {
    const url = API.Repositories.Detail.render({
      namespace: this.namespace,
      repo: repo,
    });
    return request(this.global).get<
      RepositoryEntity,
      AxiosResponse<RepositoryEntity>
    >(url);
  }

  // TODO 分页
  list(page = 0): Promise<AxiosResponse<TypeRepositories>> {
    const url = API.Repositories.List.render({
      namespace: this.namespace,
    });
    return request(this.global).get<
      TypeRepositories,
      AxiosResponse<TypeRepositories>
    >(url);
  }
}
