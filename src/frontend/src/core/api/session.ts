import { Router } from "../../config/router";

const AuthUserToken = "growerlab-token";

export interface UserInfo {
  token: string;
  namespace: string;
  email: string;
  name: string;
  public_email: string;
}

export function useSession(): Session {
  return new Session();
}

class Session {
  constructor() {
    return;
  }

  /**
   * 用户是否登录
   */
  isLogin(): boolean {
    const result = this.getCurrentUser();
    if (result === undefined) {
      return false;
    }
    return true;
  }

  /**
   * 登录，将保存token并可以设置过期时间，默认不过期
   */
  storeLogin(info: UserInfo): void {
    localStorage.setItem(AuthUserToken, JSON.stringify(info));
  }

  /**
   * 退出登录
   */
  logout(callback?: () => void) {
    localStorage.removeItem(AuthUserToken);
    if (callback === undefined) {
      callback = () => {
        location.href = Router.Home.Login;
      };
    }
    callback?.();
  }

  /**
   * 获取用户信息
   */
  getCurrentUser(): UserInfo | undefined {
    const info = localStorage.getItem(AuthUserToken);
    if (info === null) {
      console.error("Not found logined user");
      return undefined;
    }
    try {
      return JSON.parse(info) as UserInfo;
    } catch (error) {
      this.logout();
      console.error("Can't parse json for login info.");
      return undefined;
    }
  }
}
