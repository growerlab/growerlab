import { AxiosResponse } from "axios";

import { LoginInfo } from "./session";
import { API, request } from "../../api/api";
import { Notice } from "../../global/recoil/notice";

type Login = {
  email: string;
  password: string;
};

export class LoginService {
  notice: Notice;

  constructor(notice: Notice) {
    this.notice = notice;
    return;
  }

  public login(
    email: string,
    password: string
  ): Promise<AxiosResponse<LoginInfo>> {
    return request(this.notice)
      .post<Login, AxiosResponse<LoginInfo>>(API.Login, {
        email: email,
        password: password,
      })
      .then((res) => {
        return res;
      });
  }
}
