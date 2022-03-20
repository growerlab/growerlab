import { AxiosResponse } from 'axios';
import { LoginInfo } from '../../services/auth/session';
import { API, request } from '../api';

export class Login {
  private email: string;
  private password: string;

  constructor(email: string, password: string) {
    this.email = email;
    this.password = password;
  }

  do(): Promise<AxiosResponse<LoginInfo>> {
    return request()
      .post<Login, AxiosResponse<LoginInfo>>(API.Login, {
        email: this.email,
        password: this.password,
      })
      .then((res) => {
        return res;
      });
  }
}
