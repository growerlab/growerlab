import { atom, useSetRecoilState, useRecoilValue } from "recoil";

/**
 * 1. 打开网页默认(父级layout组件应初始化该选项)
 * 2. 点击menu item
 */
const userMenuSelectedState = atom<string>({
  key: "user_menu_selected",
  default: "user",
});

export const useUserMenu = () => {
  const setUserMenu = useSetRecoilState(userMenuSelectedState);
  const userMenuSelected = useRecoilValue(userMenuSelectedState);
  return { setUserMenu, userMenuSelected }
};
