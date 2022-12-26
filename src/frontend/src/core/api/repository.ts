import { AxiosResponse } from "axios";

import {
  FileEntity,
  RepositoryEntity,
  RepositoryPathTree,
  TypeRepositories,
} from "../common/types";
import { API, request } from "./api";
import { GlobalObject, useGlobal } from "../global/global";

export interface RepositoryRequest {
  namespace: string;
  name: string;
  description: string;
  public: boolean;
}

export function useRepositoryAPI(namespace: string) {
  const global = useGlobal();
  return new Repository(namespace, global);
}

class Repository {
  private readonly namespace: string;
  private readonly global: GlobalObject;

  constructor(namespace: string, global: GlobalObject) {
    this.namespace = namespace;
    this.global = global;
  }

  create(req: RepositoryRequest) {
    const url = API.Repositories.Create.render({ namespace: this.namespace });
    return request(this.global).post<RepositoryRequest>(url, req);
  }

  getDetail(repo: string) {
    const url = API.Repositories.Detail.render({
      namespace: this.namespace,
      repo: repo,
    });
    return request(this.global).get<RepositoryEntity>(url);
  }

  treeFiles(params: RepositoryPathTree) {
    const url = API.Repositories.TreeFiles.render({
      namespace: this.namespace,
      repo: params.repo,
      ref: params.ref,
      dir: params.dir,
    });
    return request(this.global).get<FileEntity[]>(url);
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
