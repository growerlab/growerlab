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
  static storeLogin(info: UserInfo): Promise<UserInfo> {
    localStorage.setItem(AuthUserToken, JSON.stringify(info));
    return Session.getUserInfo();
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
  static getUserInfo(): Promise<UserInfo> {
    const info = localStorage.getItem(AuthUserToken);

    return new Promise((resolve, reject) => {
      if (info === null) {
        return reject(new Error("not found user token"));
      }
      try {
        resolve(JSON.parse(info) as UserInfo);
      } catch (error) {
        console.warn("Can't parse json for login info.");
        reject(error);
      }
    });
  }
}
