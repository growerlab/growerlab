import { UserInfo, Session } from "../services/auth/session";
import { Notice, useNotice } from "./recoil/notice";

export type globalTypes = {
  notice?: Notice;
  getUserInfo: () => UserInfo | undefined;
};

export let global: globalTypes;

export const setup = () => {
  global = {
    getUserInfo: () => undefined,
  };
  return;
};

export function useGlobal(): globalTypes {
  global.notice = useNotice();
  global.getUserInfo = Session.getUserInfo;

  return global;
}
