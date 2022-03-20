import { Login } from "../../api/auth/login";
import { LoginInfo, Session } from "./session";
import { Message } from "../../api/common/notice";

export class LoginService {
  constructor() {
    return;
  }

  static login(email: string, password: string): Promise<LoginInfo> {
    const login = new Login(email, password);
    return login.do().then((res) => {
      const info = res.data;
      Session.storeLogin(info);
      return info;
    });
  }
}
