import { RepositoryPath, RepositoryPathGroup } from "../core/common/types";

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
    return this.r.replace(/:([^/]+)/g, (_: unknown, p: keyof T) => params[p]);
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
      Show: dynamicRouter.new<RepositoryPath>("/user/repos/:repo"),
      Branchs: dynamicRouter.new<RepositoryPath>("/user/repos/:repo/branchs"),
    },
    Project: {
      Index: "/user/projects",
    },
  },
  Namespace: {
    Repository: dynamicRouter.new<RepositoryPathGroup>("/:namespace/:repo"),
  },
};
