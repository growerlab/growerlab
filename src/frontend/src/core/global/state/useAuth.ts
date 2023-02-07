import { AxiosResponse } from "axios";

import { UserInfo } from "./index";
import { API, request, Result } from "../../api/api";
import { GlobalObject, useGlobal } from "../global";

interface RegisterArgs {
  username: string;
  email: string;
  password: string;
}

export function useAuth() {
  const global = useGlobal();
  return new UseAuth(global);
}

class UseAuth {
  constructor(private global: GlobalObject) {}

  public login(
    email: string,
    password: string
  ): Promise<AxiosResponse<UserInfo>> {
    return request(this.global).post<UseAuth, AxiosResponse<UserInfo>>(
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
