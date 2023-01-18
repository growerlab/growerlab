import {
  RepositoryCommit,
  RepositoryFile,
  RepositoryPath,
  RepositoryPathGroup,
} from "../core/common/types";
import { generatePath } from "react-router-dom";

type Params<T> = {
  [key in keyof T]: string;
};

export class dynamicRouter<T> {
  private r: string;

  constructor(r: string) {
    this.r = r;
  }

  static new<T>(r: string): dynamicRouter<T> {
    return new dynamicRouter(r);
  }

  public render(params: Params<T>) {
    return generatePath(this.r, params);
    // return this.r.replace(/:([^/]+)/g, (_: unknown, p: keyof T) => params[p]);
  }

  public string() {
    return this.r;
  }
}

export const Router = {
  Home: {
    Index: "/",
    Register: "/register",
    Login: "/login",
    ActivateUser: "/activate_user/:code",
  },
  User: {
    Index: "/user/",
    Repository: {
      Index: "/user/repos",
      New: "/user/repos/new",
      Show: dynamicRouter.new<RepositoryPath>("/user/repos/:repo"), // 默认文件树
      Tree: dynamicRouter.new<RepositoryFile>(
        "/user/repos/:repo/tree/:ref/*" // 文件树
      ),
      Blob: dynamicRouter.new<RepositoryFile>(
        "/user/repos/:repo/blob/:ref/:filepath" // 文件详情
      ),
      Commit: dynamicRouter.new<RepositoryCommit>(
        "/user/repos/:repo/commit/:commit" // commit详情
      ),
      Branches: dynamicRouter.new<RepositoryPath>("/user/repos/:repo/branches"), // 分支列表
    },
    Project: {
      Index: "/user/projects",
    },
  },
  Namespace: {
    Repository: dynamicRouter.new<RepositoryPathGroup>("/:namespace/:repo"),
  },
};
