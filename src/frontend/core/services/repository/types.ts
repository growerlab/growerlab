export interface Owner {
  name: string;
  namespace: string;
}

export interface RepositoryEntity {
  uuid: string;
  name: string;
  path: string; // repo path
  description: string;
  createdAt: number;
  public: boolean;
  // StartCount: number;
  // ForkCount: number;
  // LastUpdatedAt: number;
  pathGroup: string;
  gitHttpURL: string;
  gitSshURL: string;
  owner: Owner;
}

export interface TypeRepositories {
  repositories: RepositoryEntity[];
}

export interface TypeRepository {
  repository: RepositoryEntity;
}

export interface RepositoriesArgs {
  ownerPath: string;
}

export interface RepositoryArgs {
  ownerPath: string;
  repoPath: string;
}
