import { create } from "zustand";

/**
 * 1. 打开网页默认(父级layout组件应初始化该选项)
 * 2. 点击menu item
 */

interface globeUserMenuState {
  current: string;
  getMenuSelected: () => string;
  setMenuSelect: (current: string) => void;
}

export const useUserMenu = create<globeUserMenuState>((set, getState) => ({
  current: "user",
  getMenuSelected: () => getState().current,
  setMenuSelect: (current: string) => set((state) => ({ current: current })),
}));
