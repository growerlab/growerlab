import axios, { AxiosInstance, AxiosResponse } from "axios";
import i18n from "../i18n/i18n";
import { dynamicRouter } from "../../config/router";
import {
  RepositoriesNamespace,
  RepositoryPathGroup,
  RepositoryPathTree,
} from "../common/types";
import { GlobalObject } from "../global/global";

const baseUrl = "http://localhost:8081/api/v1";

export const API = {
  Auth: {
    Login: "/auth/login",
    Activate: "/auth/activate",
    Register: "/auth/register",
  },
  Repositories: {
    List: dynamicRouter.new<RepositoriesNamespace>(
      "/repositories/:namespace/list"
    ),
    Detail: dynamicRouter.new<RepositoryPathGroup>(
      "/repositories/:namespace/detail/:repo"
    ),
    Create: dynamicRouter.new<RepositoriesNamespace>(
      "/repositories/:namespace/create"
    ),
    TreeFiles: dynamicRouter.new<RepositoryPathTree>(
      "/repositories/:namespace/detail/:repo/tree/:ref/:dir"
    ),
  },
};

export interface Result {
  code: string;
  message: string;
}

/**
 * 封装axios的请求
 * @returns {AxiosInstance}
 */
export const request = function (global: GlobalObject): AxiosInstance {
  const { currentUser, notice } = global;
  const instance = axios.create({
    baseURL: baseUrl,
    timeout: 2000,
    timeoutErrorMessage: i18n.t("api.timeout"),
    // responseType: "json",
    headers: {
      // 'Content-Type': 'application/json',
      "Growerlab-Token": currentUser?.token || "",
    },
    validateStatus: function (status: number): boolean {
      return status >= 200 && status <= 500;
    },
  });

  instance.interceptors.response.use(
    (response: AxiosResponse) => {
      const status = response.status;

      if (status >= 300 || status < 200) {
        if (response.data === "") {
          console.error(`error: ${status}`);
          notice.error(i18n.t("message.error.ERROR"));
          return Promise.reject("network error");
        } else {
          notice.error(response.data.message);
          return Promise.reject(response);
        }
      }
      return response;
    },
    (error: any) => {
      if (error.response && error.response.data) {
        const msg = error.response.data.message;
        console.error(msg);
        notice.error(msg);
      } else {
        console.error(error.message);
        notice.error(error.message);
      }
      // 吃掉http网络错误（例如后端无法链接）
      return Promise.resolve({});
    }
  );

  return instance;
};
