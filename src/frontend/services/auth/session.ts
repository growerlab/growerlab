import { Router } from "../../config/router";
import i18n from "../../i18n/i18n";
import { NextRouter } from "next/router";
import { Message } from "../../api/common/notice";
import { reject } from "q";

const AuthUserToken = "auth-user-token";

export interface LoginInfo {
  token: string;
  namespacePath: string;
  email: string;
  name: string;
  publicEmail: string;
}

export class Session {
  constructor() {
    return;
  }

  /**
   * 用户是否登录
   * @returns {boolean}
   */
  static async isLogin(): Promise<boolean> {
    const result = await Session.getUserInfo();
    if (result === null) {
      return Promise.reject(false);
    }
    return Promise.resolve(true);
  }

  /**
   * 登录，将保存token并可以设置过期时间，默认不过期
   * @param info
   */
  static storeLogin(info: LoginInfo): Promise<LoginInfo> {
    localStorage.setItem(AuthUserToken, JSON.stringify(info));
    return Session.getUserInfo();
  }

  /**
   * 退出登录
   * @param router
   * @param callback
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
  static getUserInfo(): Promise<LoginInfo> {
    const info = localStorage.getItem(AuthUserToken);

    return new Promise((resolve, reject) => {
      if (info === null) {
        return reject(new Error("not found user token"));
      }
      try {
        resolve(JSON.parse(info) as LoginInfo);
      } catch (error) {
        console.warn("Can't parse json for login info.");
        reject(error);
      }
    });
  }
}
