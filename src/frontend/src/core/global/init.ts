import { UserInfo, useSession } from "../api/auth/session";
import { Notice, useNotice } from "./recoil/notice";

export type GlobalTypes = {
  notice: Notice;
  currentUser?: UserInfo;
};

export function useGlobal(): GlobalTypes {
  const session = useSession()
  const global: GlobalTypes = {
    notice: useNotice(),
    currentUser: session.getCurrentUser()
  };
  return global;
}
