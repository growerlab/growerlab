import { noticeState, noticeValue } from "./recoil/notice";
import { RecoilState } from "recoil";

export type globalTypes = {
  notice: RecoilState<noticeValue>;
};

export const setup = () => {
  return;
};

export const global = getGlobal();

function getGlobal(): globalTypes {
  return {
    notice: noticeState,
  };
}
