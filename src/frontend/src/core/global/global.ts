import { UserInfo, useSession } from "../api/session";
import { Notice, useNotice } from "./recoil/notice";

export type GlobalObject = {
  notice: Notice;
  currentUser?: UserInfo;
};

export function useGlobal(): GlobalObject {
  const session = useSession();
  const global: GlobalObject = {
    notice: useNotice(),
    currentUser: session.getCurrentUser(),
  };
  return global;
}
