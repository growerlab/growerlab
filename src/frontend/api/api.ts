import axios, { AxiosInstance, AxiosResponse } from "axios";
import { Message } from "./common/notice";
import i18n from "../i18n/i18n";

const baseUrl = "http://localhost:8081/api/v1/";

export const API = {
  Login: "/auth/login",
};

/**
 * 封装axios的请求
 * @returns {AxiosInstance}
 */
export const request = function (): AxiosInstance {
  const instance = axios.create({
    baseURL: baseUrl,
    timeout: 2000,
    timeoutErrorMessage: i18n.t("api.timeout"),
    // responseType: "json",
    // headers: {
    //   'Content-Type': 'application/json',
    // },
    validateStatus: function (status: number): boolean {
      return status >= 200 && status <= 500;
    },
  });

  instance.interceptors.response.use(
    (response: AxiosResponse) => {
      const status = response.status;

      if (status >= 300 || status < 200) {
        Message.Error(response.data.message);
        return Promise.reject(response);
      }
      return response;
    },
    (error: any) => {
      if (error.response && error.response.data) {
        const msg = error.response.data.message;
        Message.Error(msg);
      } else {
        Message.Error(error.message);
      }
      // 吃掉http网络错误（例如后端无法链接）
      return Promise.resolve({});
    }
  );

  return instance;
};
