import { Router } from "../../../config/router";

const AuthUserToken = "auth-user-token";

export interface UserInfo {
  token: string;
  namespace_path: string;
  email: string;
  name: string;
  public_email: string;
}

export class Session {
  constructor() {
    return;
  }

  /**
   * 用户是否登录
   */
  static async isLogin(): Promise<boolean> {
    const result = await Session.getCurrentUser();
    if (result === undefined) {
      return Promise.reject(false);
    }
    return true;
  }

  /**
   * 登录，将保存token并可以设置过期时间，默认不过期
   */
  static storeLogin(info: UserInfo): void {
    localStorage.setItem(AuthUserToken, JSON.stringify(info));
  }

  /**
   * 退出登录
   */
  static logout(callback?: () => void) {
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
  static getCurrentUser(): UserInfo | undefined {
    const info = localStorage.getItem(AuthUserToken);
    if (info === null) {
      console.error("Not found logined user");
      return undefined;
    }
    try {
      return JSON.parse(info) as UserInfo;
    } catch (error) {
      Session.logout();
      console.error("Can't parse json for login info.");
      return undefined;
    }
  }
}
