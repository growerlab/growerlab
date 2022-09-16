import { AxiosResponse } from "axios";

import { LoginInfo } from "./session";
import { API, request } from "../../api/api";
import { global } from "../../global/init";

type Login = {
  email: string;
  password: string;
};

export class LoginService {
  public login(
    email: string,
    password: string
  ): Promise<AxiosResponse<LoginInfo>> {
    return request(global.notice!)
      .post<Login, AxiosResponse<LoginInfo>>(API.Login, {
        email: email,
        password: password,
      })
      .then((res) => {
        return res;
      });
  }
}
