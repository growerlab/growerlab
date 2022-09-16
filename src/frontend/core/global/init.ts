import { Notice, useNotice } from "./recoil/notice";

type globalTypes = {
  notice: Notice | undefined;
};

export let global: globalTypes;

export const setup = () => {
  global = {
    notice: undefined,
  };
  return;
};

export function useGlobal(): globalTypes {
  global.notice = useNotice();

  return global;
}
