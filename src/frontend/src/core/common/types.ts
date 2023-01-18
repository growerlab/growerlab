export interface Owner {
  name: string;
  public_email: string;
  username: string;
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
  last_push_at: number;
  default_branch: string;
  owner: Owner;
  namespace: NamespaceEntity;
}

export interface NamespaceEntity {
  path: string; // namespace path
  type: "user" | "org";
  // owner: Owner;
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

export interface RepositoryFile extends RepositoryPath {
  ref: string; // repo ref
  "*"?: string; // repo folder path
  filepath?: string; // repo file path（only file, not folder)
}

export interface RepositoryCommit extends RepositoryPath {
  commit: string; // commit detail
}

export interface RepositoryPathGroup {
  namespace: string;
  repo: string;
}

export interface RepositoryPathTree extends RepositoryPathGroup {
  ref: string | "main";
  folder: string | "";
}

export interface FileEntity {
  id: string;
  is_file: boolean;
  last_commit_date: number;
  last_commit_hash: string;
  last_commit_message: string;
  name: string;
  tree_hash: string;
}
