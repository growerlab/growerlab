import { UserInfo, Session } from "../services/auth/session";
import { Notice, useNotice } from "./recoil/notice";

export type globalTypes = {
  notice: Notice | null;
  getUserInfo: () => Promise<UserInfo>;
};

export let global: globalTypes;

export const setup = () => {
  global = {
    notice: null,
    getUserInfo: () => Promise.reject("not found user"),
  };
  return;
};

export function useGlobal(): globalTypes {
  global.notice = useNotice();
  global.getUserInfo = Session.getUserInfo;

  return global;
}
