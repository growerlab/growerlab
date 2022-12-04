import { AxiosResponse } from "axios";

import { UserInfo } from "./session";
import { API, request, Result } from "../api";
import { GlobalObject, useGlobal } from "../../global/global";

interface RegisterArgs {
  username: string;
  email: string;
  password: string;
}

export function useAuth() {
  const global = useGlobal();
  return new Auth(global);
}

class Auth {
  constructor(private global: GlobalObject) {}

  public login(
    email: string,
    password: string
  ): Promise<AxiosResponse<UserInfo>> {
    return request(this.global).post<Auth, AxiosResponse<UserInfo>>(
      API.Auth.Login,
      {
        email: email,
        password: password,
      }
    );
  }

  public activate(code: string): Promise<AxiosResponse<Result>> {
    return request(this.global).post<Result>(API.Auth.Activate, {
      code: code,
    });
  }

  public registerUser(args: RegisterArgs): Promise<AxiosResponse<Result>> {
    return request(this.global).post<Result>(API.Auth.Register, args);
  }
}
