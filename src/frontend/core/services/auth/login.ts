import { AxiosResponse } from "axios";

import { UserInfo } from "./session";
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
  ): Promise<AxiosResponse<UserInfo>> {
    return request(global.notice!)
      .post<Login, AxiosResponse<UserInfo>>(API.Login, {
        email: email,
        password: password,
      })
      .then((res) => {
        return res;
      });
  }
}
