import { AxiosResponse } from "axios";

import { UserInfo } from "./session";
import { API, request, Result } from "../../api/api";
import { global } from "../../global/init";

interface Login {
  email: string;
  password: string;
}

interface RegisterArgs {
  username: string;
  email: string;
  password: string;
}

export class Auth {
  public login(
    email: string,
    password: string
  ): Promise<AxiosResponse<UserInfo>> {
    return request(global.notice!).post<Auth, AxiosResponse<UserInfo>>(
      API.Auth.Login,
      {
        email: email,
        password: password,
      }
    );
  }

  public activate(code: string): Promise<AxiosResponse<Result>> {
    return request(global.notice!).post<Result>(API.Auth.Activate, {
      code: code,
    });
  }

  public registerUser(args: RegisterArgs): Promise<AxiosResponse<Result>> {
    return request(global.notice!).post<Result>(API.Auth.Register, args);
  }
}
