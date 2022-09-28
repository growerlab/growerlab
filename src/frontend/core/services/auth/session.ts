import { Router } from "../../../config/router";
import { NextRouter } from "next/router";

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
    const result = await Session.getUserInfo();
    if (result === null) {
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
  static logout(router: NextRouter, callback?: () => void) {
    localStorage.removeItem(AuthUserToken);
    if (callback === undefined) {
      callback = () => {
        router.push(Router.Home.Login);
      };
    }
    callback?.();
  }

  /**
   * 获取用户信息
   */
  static getUserInfo(): UserInfo | undefined {
    const info = localStorage.getItem(AuthUserToken);

    if (info === null) {
      console.error("not found user token");
      return undefined
    }
    try {
      return JSON.parse(info) as UserInfo;
    } catch (error) {
      console.warn("Can't parse json for login info.");
      return undefined
    }
  }
}
