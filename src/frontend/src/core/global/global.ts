import { UserInfo, useSession } from "../api/session";
import { NoticeState, useNotice } from "./state/useNotice";

export type GlobalObject = {
  notice: NoticeState;
  currentUser?: UserInfo;
};

export function useGlobal(): GlobalObject {
  const session = useSession();
  const notice = useNotice();

  const global: GlobalObject = {
    notice: notice,
    currentUser: session.getCurrentUser(),
  };
  return global;
}
