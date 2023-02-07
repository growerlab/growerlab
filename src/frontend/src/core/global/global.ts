import { UserInfo, useSession, NoticeState, useNotice } from "./state";

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
