export class dynamicRouter {
  private r: string;

  constructor(r: string) {
    this.r = r;
  }

  static new(r: string): dynamicRouter {
    return new dynamicRouter(r);
  }

  public render(params: any) {
    return this.r.replace(
      /:([^/]+)/g,
      (_: any, p: string | number) => params[p]
    );
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
      Show: dynamicRouter.new("/user/repos/:repoPath"),
    },
    Project: {
      Index: "/user/projects",
    },
  },
  Namespace: {
    Repository: dynamicRouter.new("/:namespacePath/:repoPath"),
  },
};
