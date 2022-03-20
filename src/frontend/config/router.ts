class dynamicRouter {
  private r: string

  constructor(r: string) {
    this.r = r;
  }

  static new(r: string): dynamicRouter {
    return new dynamicRouter(r);
  }

  public render(params: any) {
    return this.r.replace(/:([^/]+)/g, (_: any, p: string | number) => params[p]);
  }
}


export const Router = {
  Home: {
    Index: '/',
    Register: '/register',
    Login: '/login',
    ActivateUser: '/activate_user/:code',
  },
  User: {
    Index: '/user/',
    Repository: {
      New: '/user/repos/new',
      List: '/user/repos',
      Show: '/user/repos/:repoPath',
    },
  },
  Namespace: {
    Repository: dynamicRouter.new('/:namespacePath/:repoPath'),
  },
};
