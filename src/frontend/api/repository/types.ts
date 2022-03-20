export interface Owner {
  name: string;
  namespace: string;
}

export interface RepositoryEntity {
  uuid: string;
  name: string;
  path: string;
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

export interface TypeRepositoriesArgs {
  ownerPath: string;
}

export interface RepositoryArgs {
  ownerPath?: string;
}

// export const GQL_QUERY_REPOSITORIES = gql`
//   query repositories($ownerPath: String!) {
//     repositories(ownerPath: $ownerPath) {
//       uuid
//       path
//       name
//       description
//       createdAt
//       owner {
//         name
//         username
//       }
//       pathGroup
//       gitHttpURL
//       gitSshURL
//       public
//     }
//   }
// `;
//
// export const GQL_QUERY_REPOSITORY = gql`
//   query repository($ownerPath: String!, $repoPath: String!) {
//     repository(ownerPath: $ownerPath, path: $repoPath) {
//       uuid
//       path
//       name
//       description
//       createdAt
//       owner {
//         name
//         username
//       }
//       pathGroup
//       gitHttpURL
//       gitSshURL
//       public
//     }
//   }
// `;
