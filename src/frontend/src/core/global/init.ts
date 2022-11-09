import { UserInfo, Session } from "../services/auth/session";
import { Notice, useNotice } from "./recoil/notice";

export type globalTypes = {
  notice?: Notice;
  currentUser?: UserInfo;
};

export let global: globalTypes;

export const setup = () => {
  global = {
    currentUser: undefined,
  };
  return;
};

export function useGlobal(): globalTypes {
  global.notice = useNotice();
  global.currentUser = Session.getCurrentUser();
  return global;
}
