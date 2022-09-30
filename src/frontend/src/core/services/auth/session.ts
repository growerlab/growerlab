import { Router } from "../../../config/router";
import { redirect } from "react-router-dom";

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
  static logout(callback?: () => void) {
    localStorage.removeItem(AuthUserToken);
    if (callback === undefined) {
      callback = () => {
        redirect(Router.Home.Login);
      };
    }
    callback?.();
  }

  /**
   * 获取用户信息
   */
  static getUserInfo(): UserInfo {
    const info = localStorage.getItem(AuthUserToken);

    if (info === null) {
      console.error("not found user info");
      throw new Error();
    }
    try {
      return JSON.parse(info) as UserInfo;
    } catch (error) {
      Session.logout();
      console.warn("Can't parse json for login info.");
      throw new Error();
    }
  }
}
