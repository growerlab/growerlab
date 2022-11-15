export interface Owner {
  name: string;
  namespace: string;
}

export interface RepositoryEntity {
  uuid: string;
  name: string;
  path: string; // repo path
  description: string;
  created_at: number;
  public: boolean;
  git_http_url: string;
  git_ssh_url: string;
  owner: Owner;
  namespace: NamespaceEntity;
}

export interface NamespaceEntity {
  path: string; // namespace path
  owner: Owner;
}

export interface TypeRepositories {
  repositories: RepositoryEntity[];
}

export interface TypeRepository {
  repository: RepositoryEntity;
}

export interface Namespace {
  namespace: string;
}
export type RepositoriesNamespace = Namespace;

export interface RepositoryPath {
  repo: string; // repo path
}

export interface RepositoryPathGroup {
  namespace: string;
  repo: string;
}
